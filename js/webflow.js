/*!
 * ----------------------------------------------------------------------
 * Webflow: Front-end site library
 * @license MIT
 * Other scripts may access this api using an async handler:
 *   var Webflow = Webflow || [];
 *   Webflow.push(readyFunction);
 * ----------------------------------------------------------------------
 */
var Webflow = { w: Webflow };
Webflow.init = function () {
  'use strict';

  var $ = window.$;
  var api = {};
  var modules = {};
  var primary = [];
  var secondary = this.w || [];
  var $win = $(window);
  var _ = api._ = underscore();
  var domready = false;
  var tram = window.tram;
  var Modernizr = window.Modernizr;
  tram.config.hideBackface = false;
  tram.config.keepInherited = true;

  /**
   * Webflow.define() - Define a webflow.js module
   * @param  {string} name
   * @param  {function} factory
   */
  api.define = function (name, factory) {
    var module = modules[name] = factory($, _);
    if (!module) return;
    // If running in Webflow app, subscribe to design/preview events
    if (api.env()) {
      $.isFunction(module.design) && window.addEventListener('__wf_design', module.design);
      $.isFunction(module.preview) && window.addEventListener('__wf_preview', module.preview);
    }
    // Look for a ready method on module
    if (module.ready && $.isFunction(module.ready)) {
      // If domready has already happened, call ready method
      if (domready) module.ready();
      // Otherwise push ready method into primary queue
      else primary.push(module.ready);
    }
  };

  /**
   * Webflow.require() - Load a Webflow.js module
   * @param  {string} name
   * @return {object}
   */
  api.require = function (name) {
    return modules[name];
  };

  /**
   * Webflow.push() - Add a ready handler into secondary queue
   * @param {function} ready  Callback to invoke on domready
   */
  api.push = function (ready) {
    // If domready has already happened, invoke handler
    if (domready) {
      $.isFunction(ready) && ready();
      return;
    }
    // Otherwise push into secondary queue
    secondary.push(ready);
  };

  /**
   * Webflow.env() - Get the state of the Webflow app
   * @param {string} mode [optional]
   * @return {boolean}
   */
  api.env = function (mode) {
    var designFlag = window.__wf_design;
    var inApp = typeof designFlag != 'undefined';
    if (!mode) return inApp;
    if (mode == 'design') return inApp && designFlag;
    if (mode == 'preview') return inApp && !designFlag;
    if (mode == 'slug') return inApp && window.__wf_slug;
  };

  // Feature detects + browser sniffs  ಠ_ಠ
  var userAgent = navigator.userAgent.toLowerCase();
  var appVersion = navigator.appVersion.toLowerCase();
  api.env.touch = ('ontouchstart' in window) || window.DocumentTouch && document instanceof window.DocumentTouch;
  var chrome = api.env.chrome = (window.chrome || /chrome/.test(userAgent)) && parseInt(appVersion.match(/chrome\/(\d+)\./)[1], 10);
  var ios = api.env.ios = Modernizr && Modernizr.ios;
  api.env.safari = /safari/.test(userAgent) && !chrome && !ios;

  /**
   * Webflow.script() - Append script to document head
   * @param {string} src
   */
  api.script = function (src) {
    var doc = document;
    var scriptNode = doc.createElement('script');
    scriptNode.type = 'text/javascript';
    scriptNode.src = src;
    doc.getElementsByTagName('head')[0].appendChild(scriptNode);
  };

  /**
   * Webflow.resize, Webflow.scroll - throttled event proxies
   */
  var resizeEvents = 'resize.webflow orientationchange.webflow load.webflow';
  var scrollEvents = 'scroll.webflow ' + resizeEvents;
  api.resize = eventProxy($win, resizeEvents);
  api.scroll = eventProxy($win, scrollEvents);
  api.redraw = eventProxy();

  // Create a proxy instance for throttled events
  function eventProxy(target, types) {

    // Set up throttled method (using custom frame-based _.throttle)
    var handlers = [];
    var proxy = {};
    proxy.up = _.throttle(function (evt) {
      _.each(handlers, function (h) { h(evt); });
    });

    // Bind events to target
    if (target && types) target.on(types, proxy.up);

    /**
     * Add an event handler
     * @param  {function} handler
     */
    proxy.on = function (handler) {
      if (typeof handler != 'function') return;
      if (_.contains(handlers, handler)) return;
      handlers.push(handler);
    };

    /**
     * Remove an event handler
     * @param  {function} handler
     */
    proxy.off = function (handler) {
      handlers = _.filter(handlers, function (h) {
        return h !== handler;
      });
    };
    return proxy;
  }

  // Webflow.location - Wrap window.location in api
  api.location = function (url) {
    window.location = url;
  };

  // Webflow.app - Designer-specific methods
  api.app = api.env() ? {} : null;
  if (api.app) {

    // Trigger redraw for specific elements
    var Event = window.Event;
    var redraw = new Event('__wf_redraw');
    api.app.redrawElement = function (i, el) { el.dispatchEvent(redraw); };

    // Webflow.location - Re-route location change to trigger an event
    api.location = function (url) {
      window.dispatchEvent(new CustomEvent('__wf_location', { detail: url }));
    };
  }

  // DOM ready - Call primary and secondary handlers
  $(function () {
    domready = true;
    $.each(primary.concat(secondary), function (index, value) {
      $.isFunction(value) && value();
    });
    // Trigger resize
    api.resize.up();
  });

  /*!
   * Webflow._ (aka) Underscore.js 1.5.2 (custom build)
   * _.each
   * _.map
   * _.filter
   * _.any
   * _.contains
   * _.delay
   * _.defer
   * _.throttle (webflow)
   * _.debounce
   * _.keys
   * _.has
   *
   * http://underscorejs.org
   * (c) 2009-2013 Jeremy Ashkenas, DocumentCloud and Investigative Reporters & Editors
   * Underscore may be freely distributed under the MIT license.
   */
  function underscore() {
    var _ = {};

    // Current version.
    _.VERSION = '1.5.2-Webflow';

    // Establish the object that gets returned to break out of a loop iteration.
    var breaker = {};

    // Save bytes in the minified (but not gzipped) version:
    var ArrayProto = Array.prototype, ObjProto = Object.prototype, FuncProto = Function.prototype;

    // Create quick reference variables for speed access to core prototypes.
    var
      push             = ArrayProto.push,
      slice            = ArrayProto.slice,
      concat           = ArrayProto.concat,
      toString         = ObjProto.toString,
      hasOwnProperty   = ObjProto.hasOwnProperty;

    // All **ECMAScript 5** native function implementations that we hope to use
    // are declared here.
    var
      nativeForEach      = ArrayProto.forEach,
      nativeMap          = ArrayProto.map,
      nativeReduce       = ArrayProto.reduce,
      nativeReduceRight  = ArrayProto.reduceRight,
      nativeFilter       = ArrayProto.filter,
      nativeEvery        = ArrayProto.every,
      nativeSome         = ArrayProto.some,
      nativeIndexOf      = ArrayProto.indexOf,
      nativeLastIndexOf  = ArrayProto.lastIndexOf,
      nativeIsArray      = Array.isArray,
      nativeKeys         = Object.keys,
      nativeBind         = FuncProto.bind;

    // Collection Functions
    // --------------------

    // The cornerstone, an `each` implementation, aka `forEach`.
    // Handles objects with the built-in `forEach`, arrays, and raw objects.
    // Delegates to **ECMAScript 5**'s native `forEach` if available.
    var each = _.each = _.forEach = function(obj, iterator, context) {
      /* jshint shadow:true */
      if (obj == null) return;
      if (nativeForEach && obj.forEach === nativeForEach) {
        obj.forEach(iterator, context);
      } else if (obj.length === +obj.length) {
        for (var i = 0, length = obj.length; i < length; i++) {
          if (iterator.call(context, obj[i], i, obj) === breaker) return;
        }
      } else {
        var keys = _.keys(obj);
        for (var i = 0, length = keys.length; i < length; i++) {
          if (iterator.call(context, obj[keys[i]], keys[i], obj) === breaker) return;
        }
      }
    };

    // Return the results of applying the iterator to each element.
    // Delegates to **ECMAScript 5**'s native `map` if available.
    _.map = _.collect = function(obj, iterator, context) {
      var results = [];
      if (obj == null) return results;
      if (nativeMap && obj.map === nativeMap) return obj.map(iterator, context);
      each(obj, function(value, index, list) {
        results.push(iterator.call(context, value, index, list));
      });
      return results;
    };

    // Return all the elements that pass a truth test.
    // Delegates to **ECMAScript 5**'s native `filter` if available.
    // Aliased as `select`.
    _.filter = _.select = function(obj, iterator, context) {
      var results = [];
      if (obj == null) return results;
      if (nativeFilter && obj.filter === nativeFilter) return obj.filter(iterator, context);
      each(obj, function(value, index, list) {
        if (iterator.call(context, value, index, list)) results.push(value);
      });
      return results;
    };

    // Determine if at least one element in the object matches a truth test.
    // Delegates to **ECMAScript 5**'s native `some` if available.
    // Aliased as `any`.
    var any = _.some = _.any = function(obj, iterator, context) {
      iterator || (iterator = _.identity);
      var result = false;
      if (obj == null) return result;
      if (nativeSome && obj.some === nativeSome) return obj.some(iterator, context);
      each(obj, function(value, index, list) {
        if (result || (result = iterator.call(context, value, index, list))) return breaker;
      });
      return !!result;
    };

    // Determine if the array or object contains a given value (using `===`).
    // Aliased as `include`.
    _.contains = _.include = function(obj, target) {
      if (obj == null) return false;
      if (nativeIndexOf && obj.indexOf === nativeIndexOf) return obj.indexOf(target) != -1;
      return any(obj, function(value) {
        return value === target;
      });
    };

    // Function (ahem) Functions
    // --------------------

    // Delays a function for the given number of milliseconds, and then calls
    // it with the arguments supplied.
    _.delay = function(func, wait) {
      var args = slice.call(arguments, 2);
      return setTimeout(function(){ return func.apply(null, args); }, wait);
    };

    // Defers a function, scheduling it to run after the current call stack has
    // cleared.
    _.defer = function(func) {
      return _.delay.apply(_, [func, 1].concat(slice.call(arguments, 1)));
    };

    // Returns a function, that, when invoked, will only be triggered once every
    // browser animation frame - using tram's requestAnimationFrame polyfill.
    _.throttle = function(func) {
      var wait, args, context;
      return function () {
        if (wait) return;
        wait = true;
        args = arguments;
        context = this;
        tram.frame(function () {
          wait = false;
          func.apply(context, args);
        });
      };
    };

    // Returns a function, that, as long as it continues to be invoked, will not
    // be triggered. The function will be called after it stops being called for
    // N milliseconds. If `immediate` is passed, trigger the function on the
    // leading edge, instead of the trailing.
    _.debounce = function(func, wait, immediate) {
      var timeout, args, context, timestamp, result;
      return function() {
        context = this;
        args = arguments;
        timestamp = new Date();
        var later = function() {
          var last = (new Date()) - timestamp;
          if (last < wait) {
            timeout = setTimeout(later, wait - last);
          } else {
            timeout = null;
            if (!immediate) result = func.apply(context, args);
          }
        };
        var callNow = immediate && !timeout;
        if (!timeout) {
          timeout = setTimeout(later, wait);
        }
        if (callNow) result = func.apply(context, args);
        return result;
      };
    };

    // Object Functions
    // ----------------

    // Retrieve the names of an object's properties.
    // Delegates to **ECMAScript 5**'s native `Object.keys`
    _.keys = nativeKeys || function(obj) {
      if (obj !== Object(obj)) throw new TypeError('Invalid object');
      var keys = [];
      for (var key in obj) if (_.has(obj, key)) keys.push(key);
      return keys;
    };

    // Shortcut function for checking if an object has a given property directly
    // on itself (in other words, not on a prototype).
    _.has = function(obj, key) {
      return hasOwnProperty.call(obj, key);
    };

    // Export underscore
    return _;
  }

  // Export api
  Webflow = api;
};
/*!
 * ----------------------------------------------------------------------
 * Webflow: 3rd party plugins
 */
/* jshint ignore:start */
/*!
 * tram.js v0.8.0-global
 * Cross-browser CSS3 transitions in JavaScript
 * https://github.com/bkwld/tram
 * MIT License
 */
window.tram=function(t){function i(t,i){var e=new Z.Bare;return e.init(t,i)}function e(t){return t.replace(/[A-Z]/g,function(t){return"-"+t.toLowerCase()})}function n(t){var i=parseInt(t.slice(1),16),e=255&i>>16,n=255&i>>8,r=255&i;return[e,n,r]}function r(t,i,e){return"#"+(1<<24|t<<16|i<<8|e).toString(16).slice(1)}function s(){}function a(t,i){_("Type warning: Expected: ["+t+"] Got: ["+typeof i+"] "+i)}function o(t,i,e){_("Units do not match ["+t+"]: "+i+", "+e)}function u(t,i,e){if(void 0!==i&&(e=i),void 0===t)return e;var n=e;return K.test(t)||!V.test(t)?n=parseInt(t,10):V.test(t)&&(n=1e3*parseFloat(t)),0>n&&(n=0),n===n?n:e}function c(t){for(var i=-1,e=t?t.length:0,n=[];e>++i;){var r=t[i];r&&n.push(r)}return n}var h=function(t,i,e){function n(t){return"object"==typeof t}function r(t){return"function"==typeof t}function s(){}function a(o,u){function c(){var t=new h;return r(t.init)&&t.init.apply(t,arguments),t}function h(){}u===e&&(u=o,o=Object),c.Bare=h;var l,f=s[t]=o[t],d=h[t]=c[t]=new s;return d.constructor=c,c.mixin=function(i){return h[t]=c[t]=a(c,i)[t],c},c.open=function(t){if(l={},r(t)?l=t.call(c,d,f,c,o):n(t)&&(l=t),n(l))for(var e in l)i.call(l,e)&&(d[e]=l[e]);return r(d.init)||(d.init=o),c},c.open(u)}return a}("prototype",{}.hasOwnProperty),l={ease:["ease",function(t,i,e,n){var r=(t/=n)*t,s=r*t;return i+e*(-2.75*s*r+11*r*r+-15.5*s+8*r+.25*t)}],"ease-in":["ease-in",function(t,i,e,n){var r=(t/=n)*t,s=r*t;return i+e*(-1*s*r+3*r*r+-3*s+2*r)}],"ease-out":["ease-out",function(t,i,e,n){var r=(t/=n)*t,s=r*t;return i+e*(.3*s*r+-1.6*r*r+2.2*s+-1.8*r+1.9*t)}],"ease-in-out":["ease-in-out",function(t,i,e,n){var r=(t/=n)*t,s=r*t;return i+e*(2*s*r+-5*r*r+2*s+2*r)}],linear:["linear",function(t,i,e,n){return e*t/n+i}],"ease-in-quad":["cubic-bezier(0.550, 0.085, 0.680, 0.530)",function(t,i,e,n){return e*(t/=n)*t+i}],"ease-out-quad":["cubic-bezier(0.250, 0.460, 0.450, 0.940)",function(t,i,e,n){return-e*(t/=n)*(t-2)+i}],"ease-in-out-quad":["cubic-bezier(0.455, 0.030, 0.515, 0.955)",function(t,i,e,n){return 1>(t/=n/2)?e/2*t*t+i:-e/2*(--t*(t-2)-1)+i}],"ease-in-cubic":["cubic-bezier(0.550, 0.055, 0.675, 0.190)",function(t,i,e,n){return e*(t/=n)*t*t+i}],"ease-out-cubic":["cubic-bezier(0.215, 0.610, 0.355, 1)",function(t,i,e,n){return e*((t=t/n-1)*t*t+1)+i}],"ease-in-out-cubic":["cubic-bezier(0.645, 0.045, 0.355, 1)",function(t,i,e,n){return 1>(t/=n/2)?e/2*t*t*t+i:e/2*((t-=2)*t*t+2)+i}],"ease-in-quart":["cubic-bezier(0.895, 0.030, 0.685, 0.220)",function(t,i,e,n){return e*(t/=n)*t*t*t+i}],"ease-out-quart":["cubic-bezier(0.165, 0.840, 0.440, 1)",function(t,i,e,n){return-e*((t=t/n-1)*t*t*t-1)+i}],"ease-in-out-quart":["cubic-bezier(0.770, 0, 0.175, 1)",function(t,i,e,n){return 1>(t/=n/2)?e/2*t*t*t*t+i:-e/2*((t-=2)*t*t*t-2)+i}],"ease-in-quint":["cubic-bezier(0.755, 0.050, 0.855, 0.060)",function(t,i,e,n){return e*(t/=n)*t*t*t*t+i}],"ease-out-quint":["cubic-bezier(0.230, 1, 0.320, 1)",function(t,i,e,n){return e*((t=t/n-1)*t*t*t*t+1)+i}],"ease-in-out-quint":["cubic-bezier(0.860, 0, 0.070, 1)",function(t,i,e,n){return 1>(t/=n/2)?e/2*t*t*t*t*t+i:e/2*((t-=2)*t*t*t*t+2)+i}],"ease-in-sine":["cubic-bezier(0.470, 0, 0.745, 0.715)",function(t,i,e,n){return-e*Math.cos(t/n*(Math.PI/2))+e+i}],"ease-out-sine":["cubic-bezier(0.390, 0.575, 0.565, 1)",function(t,i,e,n){return e*Math.sin(t/n*(Math.PI/2))+i}],"ease-in-out-sine":["cubic-bezier(0.445, 0.050, 0.550, 0.950)",function(t,i,e,n){return-e/2*(Math.cos(Math.PI*t/n)-1)+i}],"ease-in-expo":["cubic-bezier(0.950, 0.050, 0.795, 0.035)",function(t,i,e,n){return 0===t?i:e*Math.pow(2,10*(t/n-1))+i}],"ease-out-expo":["cubic-bezier(0.190, 1, 0.220, 1)",function(t,i,e,n){return t===n?i+e:e*(-Math.pow(2,-10*t/n)+1)+i}],"ease-in-out-expo":["cubic-bezier(1, 0, 0, 1)",function(t,i,e,n){return 0===t?i:t===n?i+e:1>(t/=n/2)?e/2*Math.pow(2,10*(t-1))+i:e/2*(-Math.pow(2,-10*--t)+2)+i}],"ease-in-circ":["cubic-bezier(0.600, 0.040, 0.980, 0.335)",function(t,i,e,n){return-e*(Math.sqrt(1-(t/=n)*t)-1)+i}],"ease-out-circ":["cubic-bezier(0.075, 0.820, 0.165, 1)",function(t,i,e,n){return e*Math.sqrt(1-(t=t/n-1)*t)+i}],"ease-in-out-circ":["cubic-bezier(0.785, 0.135, 0.150, 0.860)",function(t,i,e,n){return 1>(t/=n/2)?-e/2*(Math.sqrt(1-t*t)-1)+i:e/2*(Math.sqrt(1-(t-=2)*t)+1)+i}],"ease-in-back":["cubic-bezier(0.600, -0.280, 0.735, 0.045)",function(t,i,e,n,r){return void 0===r&&(r=1.70158),e*(t/=n)*t*((r+1)*t-r)+i}],"ease-out-back":["cubic-bezier(0.175, 0.885, 0.320, 1.275)",function(t,i,e,n,r){return void 0===r&&(r=1.70158),e*((t=t/n-1)*t*((r+1)*t+r)+1)+i}],"ease-in-out-back":["cubic-bezier(0.680, -0.550, 0.265, 1.550)",function(t,i,e,n,r){return void 0===r&&(r=1.70158),1>(t/=n/2)?e/2*t*t*(((r*=1.525)+1)*t-r)+i:e/2*((t-=2)*t*(((r*=1.525)+1)*t+r)+2)+i}]},f={"ease-in-back":"cubic-bezier(0.600, 0, 0.735, 0.045)","ease-out-back":"cubic-bezier(0.175, 0.885, 0.320, 1)","ease-in-out-back":"cubic-bezier(0.680, 0, 0.265, 1)"},d=document,p=window,b="bkwld-tram",m=/[\-\.0-9]/g,v=/[A-Z]/,g="number",y=/^(rgb|#)/,w=/(em|cm|mm|in|pt|pc|px)$/,k=/(em|cm|mm|in|pt|pc|px|%)$/,x=/(deg|rad|turn)$/,z="unitless",q=/(all|none) 0s ease 0s/,$=/^(width|height)$/,M=" ",A=d.createElement("a"),B=["Webkit","Moz","O","ms"],R=["-webkit-","-moz-","-o-","-ms-"],F=function(t){if(t in A.style)return{dom:t,css:t};var i,e,n="",r=t.split("-");for(i=0;r.length>i;i++)n+=r[i].charAt(0).toUpperCase()+r[i].slice(1);for(i=0;B.length>i;i++)if(e=B[i]+n,e in A.style)return{dom:e,css:R[i]+t}},S=i.support={bind:Function.prototype.bind,transform:F("transform"),transition:F("transition"),backface:F("backface-visibility"),timing:F("transition-timing-function")};if(S.transition){var j=S.timing.dom;if(A.style[j]=l["ease-in-back"][0],!A.style[j])for(var I in f)l[I][0]=f[I]}var G=i.frame=function(){var t=p.requestAnimationFrame||p.webkitRequestAnimationFrame||p.mozRequestAnimationFrame||p.oRequestAnimationFrame||p.msRequestAnimationFrame;return t&&S.bind?t.bind(p):function(t){p.setTimeout(t,16)}}(),T=i.now=function(){var t=p.performance,i=t&&(t.now||t.webkitNow||t.msNow||t.mozNow);return i&&S.bind?i.bind(t):Date.now||function(){return+new Date}}(),U=h(function(i){function n(t,i){var e=c((""+t).split(M)),n=e[0];i=i||{};var r=W[n];if(!r)return _("Unsupported property: "+n);if(!i.weak||!this.props[n]){var s=r[0],a=this.props[n];return a||(a=this.props[n]=new s.Bare),a.init(this.$el,e,r,i),a}}function r(t,i,e){if(t){var r=typeof t;if(i||(this.timer&&this.timer.destroy(),this.queue=[],this.active=!1),"number"==r&&i)return this.timer=new Y({duration:t,context:this,complete:o}),this.active=!0,void 0;if("string"==r&&i){switch(t){case"hide":d.call(this);break;case"stop":h.call(this);break;case"redraw":p.call(this);break;default:n.call(this,t,e&&e[1])}return o.call(this)}if("function"==r)return t.call(this,this),void 0;if("object"==r){var s=0;m.call(this,t,function(t,i){t.span>s&&(s=t.span),t.stop(),t.animate(i)},function(t){"wait"in t&&(s=u(t.wait,0))}),b.call(this),s>0&&(this.timer=new Y({duration:s,context:this}),this.active=!0,i&&(this.timer.complete=o));var a=this,c=!1,l={};G(function(){m.call(a,t,function(t){t.active&&(c=!0,l[t.name]=t.nextStyle)}),c&&a.$el.css(l)})}}}function s(t){t=u(t,0),this.active?this.queue.push({options:t}):(this.timer=new Y({duration:t,context:this,complete:o}),this.active=!0)}function a(t){return this.active?(this.queue.push({options:t,args:arguments}),this.timer.complete=o,void 0):_("No active transition timer. Use start() or wait() before then().")}function o(){if(this.timer&&this.timer.destroy(),this.active=!1,this.queue.length){var t=this.queue.shift();r.call(this,t.options,!0,t.args)}}function h(t){this.timer&&this.timer.destroy(),this.queue=[],this.active=!1;var i;"string"==typeof t?(i={},i[t]=1):i="object"==typeof t&&null!=t?t:this.props,m.call(this,i,g),b.call(this)}function l(t){h.call(this,t),m.call(this,t,y,w)}function f(t){"string"!=typeof t&&(t="block"),this.el.style.display=t}function d(){h.call(this),this.el.style.display="none"}function p(){this.el.offsetHeight}function b(){var t,i,e=[];this.upstream&&e.push(this.upstream);for(t in this.props)i=this.props[t],i.active&&e.push(i.string);e=e.join(","),this.style!==e&&(this.style=e,this.el.style[S.transition.dom]=e)}function m(t,i,r){var s,a,o,u,c=i!==g,h={};for(s in t)o=t[s],s in J?(h.transform||(h.transform={}),h.transform[s]=o):(v.test(s)&&(s=e(s)),s in W?h[s]=o:(u||(u={}),u[s]=o));for(s in h){if(o=h[s],a=this.props[s],!a){if(!c)continue;a=n.call(this,s)}i.call(this,a,o)}r&&u&&r.call(this,u)}function g(t){t.stop()}function y(t,i){t.set(i)}function w(t){this.$el.css(t)}function k(t,e){i[t]=function(){return this.children?x.call(this,e,arguments):(this.el&&e.apply(this,arguments),this)}}function x(t,i){var e,n=this.children.length;for(e=0;n>e;e++)t.apply(this.children[e],i);return this}i.init=function(i){if(this.$el=t(i),this.el=this.$el[0],this.props={},this.queue=[],this.style="",this.active=!1,C.keepInherited&&!C.fallback){var e=L(this.el,"transition");e&&!q.test(e)&&(this.upstream=e)}S.backface&&C.hideBackface&&D(this.el,S.backface.css,"hidden")},k("add",n),k("start",r),k("wait",s),k("then",a),k("next",o),k("stop",h),k("set",l),k("show",f),k("hide",d),k("redraw",p)}),Z=h(U,function(i){function e(i,e){var n=t.data(i,b)||t.data(i,b,new U.Bare);return n.el||n.init(i),e?n.start(e):n}i.init=function(i,n){var r=t(i);if(!r.length)return this;if(1===r.length)return e(r[0],n);var s=[];return r.each(function(t,i){s.push(e(i,n))}),this.children=s,this}}),H=h(function(t){function i(){var t=this.get();this.update("auto");var i=this.get();return this.update(t),i}function e(t,i,e){return void 0!==i&&(e=i),t in l?t:e}function n(t){var i=/rgba?\((\d+),\s*(\d+),\s*(\d+)/.exec(t);return(i?r(i[1],i[2],i[3]):t).replace(/#(\w)(\w)(\w)$/,"#$1$1$2$2$3$3")}var s={duration:500,ease:"ease",delay:0};t.init=function(t,i,n,r){this.$el=t,this.el=t[0];var a=i[0];n[2]&&(a=n[2]),Q[a]&&(a=Q[a]),this.name=a,this.type=n[1],this.duration=u(i[1],this.duration,s.duration),this.ease=e(i[2],this.ease,s.ease),this.delay=u(i[3],this.delay,s.delay),this.span=this.duration+this.delay,this.active=!1,this.nextStyle=null,this.auto=$.test(this.name),this.unit=r.unit||this.unit||C.defaultUnit,this.angle=r.angle||this.angle||C.defaultAngle,C.fallback||r.fallback?this.animate=this.fallback:(this.animate=this.transition,this.string=this.name+M+this.duration+"ms"+("ease"!=this.ease?M+l[this.ease][0]:"")+(this.delay?M+this.delay+"ms":""))},t.set=function(t){t=this.convert(t,this.type),this.update(t),this.redraw()},t.transition=function(t){this.active=!0,t=this.convert(t,this.type),this.auto&&("auto"==this.el.style[this.name]&&(this.update(this.get()),this.redraw()),"auto"==t&&(t=i.call(this))),this.nextStyle=t},t.fallback=function(t){var e=this.el.style[this.name]||this.convert(this.get(),this.type);t=this.convert(t,this.type),this.auto&&("auto"==e&&(e=this.convert(this.get(),this.type)),"auto"==t&&(t=i.call(this))),this.tween=new X({from:e,to:t,duration:this.duration,delay:this.delay,ease:this.ease,update:this.update,context:this})},t.get=function(){return L(this.el,this.name)},t.update=function(t){D(this.el,this.name,t)},t.stop=function(){(this.active||this.nextStyle)&&(this.active=!1,this.nextStyle=null,D(this.el,this.name,this.get()));var t=this.tween;t&&t.context&&t.destroy()},t.convert=function(t,i){if("auto"==t&&this.auto)return t;var e,r="number"==typeof t,s="string"==typeof t;switch(i){case g:if(r)return t;if(s&&""===t.replace(m,""))return+t;e="number(unitless)";break;case y:if(s){if(""===t&&this.original)return this.original;if(i.test(t))return"#"==t.charAt(0)&&7==t.length?t:n(t)}e="hex or rgb string";break;case w:if(r)return t+this.unit;if(s&&i.test(t))return t;e="number(px) or string(unit)";break;case k:if(r)return t+this.unit;if(s&&i.test(t))return t;e="number(px) or string(unit or %)";break;case x:if(r)return t+this.angle;if(s&&i.test(t))return t;e="number(deg) or string(angle)";break;case z:if(r)return t;if(s&&k.test(t))return t;e="number(unitless) or string(unit or %)"}return a(e,t),t},t.redraw=function(){this.el.offsetHeight}}),N=h(H,function(t,i){t.init=function(){i.init.apply(this,arguments),this.original||(this.original=this.convert(this.get(),y))}}),O=h(H,function(t,i){t.init=function(){i.init.apply(this,arguments),this.animate=this.fallback},t.get=function(){return this.$el[this.name]()},t.update=function(t){this.$el[this.name](t)}}),P=h(H,function(t,i){function e(t,i){var e,n,r,s,a;for(e in t)s=J[e],r=s[0],n=s[1]||e,a=this.convert(t[e],r),i.call(this,n,a,r)}t.init=function(){i.init.apply(this,arguments),this.current||(this.current={},J.perspective&&C.perspective&&(this.current.perspective=C.perspective,D(this.el,this.name,this.style(this.current)),this.redraw()))},t.set=function(t){e.call(this,t,function(t,i){this.current[t]=i}),D(this.el,this.name,this.style(this.current)),this.redraw()},t.transition=function(t){var i=this.values(t);this.tween=new E({current:this.current,values:i,duration:this.duration,delay:this.delay,ease:this.ease});var e,n={};for(e in this.current)n[e]=e in i?i[e]:this.current[e];this.active=!0,this.nextStyle=this.style(n)},t.fallback=function(t){var i=this.values(t);this.tween=new E({current:this.current,values:i,duration:this.duration,delay:this.delay,ease:this.ease,update:this.update,context:this})},t.update=function(){D(this.el,this.name,this.style(this.current))},t.style=function(t){var i,e="";for(i in t)e+=i+"("+t[i]+") ";return e},t.values=function(t){var i,n={};return e.call(this,t,function(t,e,r){n[t]=e,void 0===this.current[t]&&(i=0,~t.indexOf("scale")&&(i=1),this.current[t]=this.convert(i,r))}),n}}),X=h(function(i){function e(t){1===d.push(t)&&G(a)}function a(){var t,i,e,n=d.length;if(n)for(G(a),i=T(),t=n;t--;)e=d[t],e&&e.render(i)}function u(i){var e,n=t.inArray(i,d);n>=0&&(e=d.slice(n+1),d.length=n,e.length&&(d=d.concat(e)))}function c(t){return Math.round(t*p)/p}function h(t,i,e){return r(t[0]+e*(i[0]-t[0]),t[1]+e*(i[1]-t[1]),t[2]+e*(i[2]-t[2]))}var f={ease:l.ease[1],from:0,to:1};i.init=function(t){this.duration=t.duration||0,this.delay=t.delay||0;var i=t.ease||f.ease;l[i]&&(i=l[i][1]),"function"!=typeof i&&(i=f.ease),this.ease=i,this.update=t.update||s,this.complete=t.complete||s,this.context=t.context||this,this.name=t.name;var e=t.from,n=t.to;void 0===e&&(e=f.from),void 0===n&&(n=f.to),this.unit=t.unit||"","number"==typeof e&&"number"==typeof n?(this.begin=e,this.change=n-e):this.format(n,e),this.value=this.begin+this.unit,this.start=T(),t.autoplay!==!1&&this.play()},i.play=function(){this.active||(this.start||(this.start=T()),this.active=!0,e(this))},i.stop=function(){this.active&&(this.active=!1,u(this))},i.render=function(t){var i,e=t-this.start;if(this.delay){if(this.delay>=e)return;e-=this.delay}if(this.duration>e){var n=this.ease(e,0,1,this.duration);return i=this.startRGB?h(this.startRGB,this.endRGB,n):c(this.begin+n*this.change),this.value=i+this.unit,this.update.call(this.context,this.value),void 0}i=this.endHex||this.begin+this.change,this.value=i+this.unit,this.update.call(this.context,this.value),this.complete.call(this.context),this.destroy()},i.format=function(t,i){if(i+="",t+="","#"==t.charAt(0))return this.startRGB=n(i),this.endRGB=n(t),this.endHex=t,this.begin=0,this.change=1,void 0;if(!this.unit){var e=i.replace(m,""),r=t.replace(m,"");e!==r&&o("tween",i,t),this.unit=e}i=parseFloat(i),t=parseFloat(t),this.begin=this.value=i,this.change=t-i},i.destroy=function(){this.stop(),this.context=null,this.ease=this.update=this.complete=s};var d=[],p=1e3}),Y=h(X,function(t){t.init=function(t){this.duration=t.duration||0,this.complete=t.complete||s,this.context=t.context,this.play()},t.render=function(t){var i=t-this.start;this.duration>i||(this.complete.call(this.context),this.destroy())}}),E=h(X,function(t,i){t.init=function(t){this.context=t.context,this.update=t.update,this.tweens=[],this.current=t.current;var i,e;for(i in t.values)e=t.values[i],this.current[i]!==e&&this.tweens.push(new X({name:i,from:this.current[i],to:e,duration:t.duration,delay:t.delay,ease:t.ease,autoplay:!1}));this.play()},t.render=function(t){var i,e,n=this.tweens.length,r=!1;for(i=n;i--;)e=this.tweens[i],e.context&&(e.render(t),this.current[e.name]=e.value,r=!0);return r?(this.update&&this.update.call(this.context),void 0):this.destroy()},t.destroy=function(){if(i.destroy.call(this),this.tweens){var t,e=this.tweens.length;for(t=e;t--;)this.tweens[t].destroy();this.tweens=null,this.current=null}}}),C=i.config={defaultUnit:"px",defaultAngle:"deg",keepInherited:!1,hideBackface:!1,perspective:"",fallback:!S.transition,agentTests:[]};i.fallback=function(t){if(!S.transition)return C.fallback=!0;C.agentTests.push("("+t+")");var i=RegExp(C.agentTests.join("|"),"i");C.fallback=i.test(navigator.userAgent)},i.fallback("6.0.[2-5] Safari"),i.tween=function(t){return new X(t)},i.delay=function(t,i,e){return new Y({complete:i,duration:t,context:e})},t.fn.tram=function(t){return i.call(null,this,t)};var D=t.style,L=t.css,Q={transform:S.transform&&S.transform.css},W={color:[N,y],background:[N,y,"background-color"],"outline-color":[N,y],"border-color":[N,y],"border-top-color":[N,y],"border-right-color":[N,y],"border-bottom-color":[N,y],"border-left-color":[N,y],"border-width":[H,w],"border-top-width":[H,w],"border-right-width":[H,w],"border-bottom-width":[H,w],"border-left-width":[H,w],"border-spacing":[H,w],"letter-spacing":[H,w],margin:[H,w],"margin-top":[H,w],"margin-right":[H,w],"margin-bottom":[H,w],"margin-left":[H,w],padding:[H,w],"padding-top":[H,w],"padding-right":[H,w],"padding-bottom":[H,w],"padding-left":[H,w],"outline-width":[H,w],opacity:[H,g],top:[H,k],right:[H,k],bottom:[H,k],left:[H,k],"font-size":[H,k],"text-indent":[H,k],"word-spacing":[H,k],width:[H,k],"min-width":[H,k],"max-width":[H,k],height:[H,k],"min-height":[H,k],"max-height":[H,k],"line-height":[H,z],"scroll-top":[O,g,"scrollTop"],"scroll-left":[O,g,"scrollLeft"]},J={};S.transform&&(W.transform=[P],J={x:[k,"translateX"],y:[k,"translateY"],rotate:[x],rotateX:[x],rotateY:[x],scale:[g],scaleX:[g],scaleY:[g],skew:[x],skewX:[x],skewY:[x]}),S.transform&&S.backface&&(J.z=[k,"translateZ"],J.rotateZ=[x],J.scaleZ=[g],J.perspective=[w]);var K=/ms/,V=/s|\./,_=function(){var t="warn",i=window.console;return i&&i[t]?function(e){i[t](e)}:s}();return t.tram=i}(window.jQuery);/*!
 * jQuery-ajaxTransport-XDomainRequest - v1.0.1 - 2013-10-17
 * https://github.com/MoonScript/jQuery-ajaxTransport-XDomainRequest
 * Copyright (c) 2013 Jason Moon (@JSONMOON)
 * Licensed MIT (/blob/master/LICENSE.txt)
 */
(function($){if(!$.support.cors&&$.ajaxTransport&&window.XDomainRequest){var n=/^https?:\/\//i;var o=/^get|post$/i;var p=new RegExp('^'+location.protocol,'i');var q=/text\/html/i;var r=/\/json/i;var s=/\/xml/i;$.ajaxTransport('* text html xml json',function(i,j,k){if(i.crossDomain&&i.async&&o.test(i.type)&&n.test(i.url)&&p.test(i.url)){var l=null;var m=(j.dataType||'').toLowerCase();return{send:function(f,g){l=new XDomainRequest();if(/^\d+$/.test(j.timeout)){l.timeout=j.timeout}l.ontimeout=function(){g(500,'timeout')};l.onload=function(){var a='Content-Length: '+l.responseText.length+'\r\nContent-Type: '+l.contentType;var b={code:200,message:'success'};var c={text:l.responseText};try{if(m==='html'||q.test(l.contentType)){c.html=l.responseText}else if(m==='json'||(m!=='text'&&r.test(l.contentType))){try{c.json=$.parseJSON(l.responseText)}catch(e){b.code=500;b.message='parseerror'}}else if(m==='xml'||(m!=='text'&&s.test(l.contentType))){var d=new ActiveXObject('Microsoft.XMLDOM');d.async=false;try{d.loadXML(l.responseText)}catch(e){d=undefined}if(!d||!d.documentElement||d.getElementsByTagName('parsererror').length){b.code=500;b.message='parseerror';throw'Invalid XML: '+l.responseText;}c.xml=d}}catch(parseMessage){throw parseMessage;}finally{g(b.code,b.message,c,a)}};l.onprogress=function(){};l.onerror=function(){g(500,'error',{text:l.responseText})};var h='';if(j.data){h=($.type(j.data)==='string')?j.data:$.param(j.data)}l.open(i.type,i.url);l.send(h)},abort:function(){if(l){l.abort()}}}}})}})(jQuery);
/* jshint ignore:end */
/**
 * ----------------------------------------------------------------------
 * Init lib after plugins
 */
Webflow.init();
/**
 * ----------------------------------------------------------------------
 * Webflow: Interactions
 */
Webflow.define('ix', function ($, _) {
  'use strict';

  var api = {};
  var designer;
  var $win = $(window);
  var namespace = '.w-ix';
  var tram = window.tram;
  var env = Webflow.env;
  var ios = env.ios;
  var inApp = env();
  var emptyFix = env.chrome && env.chrome < 35;
  var transNone = 'none 0s ease 0s';
  var introEvent = 'w-ix-intro' + namespace;
  var outroEvent = 'w-ix-outro' + namespace;
  var eventQueue = [];
  var $subs = $();
  var config = {};
  var anchors = [];
  var loads = [];
  var readys = [];
  var unique = 0;
  var store;

  // -----------------------------------
  // Module methods

  api.init = function (list) {
    setTimeout(function () { init(list); }, 1);
  };

  api.preview = function () {
    designer = false;
    setTimeout(function () { init(window.__wf_ix); }, 1);
  };

  api.design = function () {
    designer = true;
    $subs.each(teardown);
    Webflow.scroll.off(scroll);
    asyncEvents();
    anchors = [];
    loads = [];
    readys = [];
  };

  api.run = run;
  api.events = {};
  api.style = inApp ? styleApp : stylePub;

  // -----------------------------------
  // Private methods

  function init(list) {
    if (!list) return;
    store = {};

    // Map all interactions to a hash using slug as key.
    config = {};
    _.each(list, function (item) {
      config[item.slug] = item.value;
    });

    // Build each element's interaction keying from data attribute
    var els = $('[data-ix]');
    els.each(teardown);
    els.each(build);

    // Listen for scroll events if any anchors exist
    if (anchors.length) {
      Webflow.scroll.on(scroll);
      setTimeout(scroll, 1);
    }

    // Handle loads or readys if they exist
    if (loads.length) $win.on('load', runLoads);
    if (readys.length) setTimeout(runReadys, 1);

    // Init module events
    initEvents();
  }

  function build(i, el) {
    var $el = $(el);
    var id = $el.attr('data-ix');
    var ix = config[id];
    if (!ix) return;
    var triggers = ix.triggers;
    if (!triggers) return;
    var state = store[id] || (store[id] = {});
    var $proxy;

    // Set initial styles, unless we detect an iOS device + any non-iOS triggers
    var setStyles = !(ios && _.any(triggers, isNonIOS));
    if (setStyles) api.style($el, ix.style);

    _.each(triggers, function (trigger) {
      var type = trigger.type;
      var stepsB = trigger.stepsB && trigger.stepsB.length;

      function runA() { run(trigger, $el, { group: 'A' }); }
      function runB() { run(trigger, $el, { group: 'B' }); }

      if (type == 'load') {
        (trigger.preload && !inApp) ? loads.push(runA) : readys.push(runA);
        return;
      }

      if (type == 'click') {
        var stateKey = 'click:' + unique++;
        if (trigger.descend) stateKey += ':descend';
        if (trigger.siblings) stateKey += ':siblings';
        if (trigger.selector) stateKey += ':' + trigger.selector;

        $el.on('click' + namespace, function (evt) {
          if ($el.attr('href') === '#') evt.preventDefault();

          run(trigger, $el, { group: state[stateKey] ? 'B' : 'A' });
          if (stepsB) state[stateKey] = !state[stateKey];
        });
        $subs = $subs.add($el);
        return;
      }

      if (type == 'hover') {
        $el.on('mouseenter' + namespace, runA);
        $el.on('mouseleave' + namespace, runB);
        $subs = $subs.add($el);
        return;
      }

      if (type == 'tabs') {
        $proxy = $el.closest('.w-tab-link, .w-tab-pane');
        $proxy.on(introEvent, runA).on(outroEvent, runB);
        $subs = $subs.add($proxy);
        return;
      }

      if (type == 'slider') {
        $proxy = $el.closest('.w-slide');
        $proxy.on(introEvent, runA).on(outroEvent, runB);
        $subs = $subs.add($proxy);
        return;
      }

      // Ignore the following triggers on iOS devices
      if (ios) return;

      if (type == 'scroll') {
        anchors.push({
          el: $el, trigger: trigger, state: { active: false },
          offsetTop: convert(trigger.offsetTop),
          offsetBot: convert(trigger.offsetBot)
        });
        return;
      }
    });
  }

  function isNonIOS(trigger) {
    return trigger.type == 'scroll';
  }

  function convert(offset) {
    if (!offset) return 0;
    offset = offset + '';
    var result = parseInt(offset, 10);
    if (result !== result) return 0;
    if (offset.indexOf('%') > 0) {
      result = result / 100;
      if (result >= 1) result = 0.999;
    }
    return result;
  }

  function teardown(i, el) {
    $(el).off(namespace);
  }

  function scroll() {
    var viewTop = $win.scrollTop();
    var viewHeight = $win.height();

    // Check each anchor for a valid scroll trigger
    var count = anchors.length;
    for (var i = 0; i < count; i++) {
      var anchor = anchors[i];
      var $el = anchor.el;
      var trigger = anchor.trigger;
      var stepsB = trigger.stepsB && trigger.stepsB.length;
      var state = anchor.state;
      var top = $el.offset().top;
      var height = $el.outerHeight();
      var offsetTop = anchor.offsetTop;
      var offsetBot = anchor.offsetBot;
      if (offsetTop < 1 && offsetTop > 0) offsetTop *= viewHeight;
      if (offsetBot < 1 && offsetBot > 0) offsetBot *= viewHeight;
      var active = (top + height - offsetTop >= viewTop && top + offsetBot <= viewTop + viewHeight);
      if (active === state.active) continue;
      if (active === false && !stepsB) continue;
      state.active = active;
      run(trigger, $el, { group: active ? 'A' : 'B' });
    }
  }

  function runLoads() {
    var count = loads.length;
    for (var i = 0; i < count; i++) {
      loads[i]();
    }
  }

  function runReadys() {
    var count = readys.length;
    for (var i = 0; i < count; i++) {
      readys[i]();
    }
  }

  function run(trigger, $el, opts, replay) {
    opts = opts || {};
    var done = opts.done;

    // Do not run in designer unless forced
    if (designer && !opts.force) return;

    // Operate on a set of grouped steps
    var group = opts.group || 'A';
    var loop = trigger['loop' + group];
    var steps = trigger['steps' + group];
    if (!steps || !steps.length) return;
    if (steps.length < 2) loop = false;

    // One-time init before any loops
    if (!replay) {

      // Find selector within element descendants, siblings, or query whole document
      var selector = trigger.selector;
      if (selector) {
        $el = (
          trigger.descend ? $el.find(selector) :
          trigger.siblings ? $el.siblings(selector) :
          $(selector)
        );
        if (inApp) $el.attr('data-ix-affect', 1);
      }

      // Apply empty fix for certain Chrome versions
      if (emptyFix) $el.addClass('w-ix-emptyfix');
    }

    var _tram = tram($el);

    // Add steps
    var meta = {};
    for (var i = 0; i < steps.length; i++) {
      addStep(_tram, steps[i], meta);
    }

    function fin() {
      // Run trigger again if looped
      if (loop) return run(trigger, $el, opts, true);

      // Reset any 'auto' values
      if (meta.width == 'auto') _tram.set({ width: 'auto' });
      if (meta.height == 'auto') _tram.set({ height: 'auto' });

      // Run callback
      done && done();
    }

    // Add final step to queue if tram has started
    meta.start ? _tram.then(fin) : fin();
  }

  function addStep(_tram, step, meta) {
    var addMethod = 'add';
    var startMethod = 'start';

    // Once the transition has started, we will always use then() to add to the queue.
    if (meta.start) addMethod = startMethod = 'then';

    // Parse transitions string on the current step
    var transitions = step.transition;
    if (transitions) {
      transitions = transitions.split(',');
      for (var i = 0; i < transitions.length; i++) {
        _tram[addMethod](transitions[i]);
      }
    }

    // Build a clean object to pass to the tram method
    var clean = tramify(step) || {};

    // Store last width and height values
    if (clean.width != null) meta.width = clean.width;
    if (clean.height != null) meta.height = clean.height;

    // When transitions are not present, set values immediately and continue queue.
    if (transitions == null) {

      // If we have started, wrap set() in then() and reset queue
      if (meta.start) {
        _tram.then(function () {
          var queue = this.queue;
          this.set(clean);
          if (clean.display) {
            _tram.redraw();
            Webflow.redraw.up();
          }
          this.queue = queue;
          this.next();
        });
      } else {
        _tram.set(clean);

        // Always redraw after setting display
        if (clean.display) {
          _tram.redraw();
          Webflow.redraw.up();
        }
      }

      // Use the wait() method to kick off queue in absence of transitions.
      var wait = clean.wait;
      if (wait != null) {
        _tram.wait(wait);
        meta.start = true;
      }

    // Otherwise, when transitions are present
    } else {

      // If display is present, handle it separately
      if (clean.display) {
        var display = clean.display;
        delete clean.display;

        // If we've already started, we need to wrap it in a then()
        if (meta.start) {
          _tram.then(function () {
            var queue = this.queue;
            this.set({ display: display }).redraw();
            Webflow.redraw.up();
            this.queue = queue;
            this.next();
          });
        } else {
          _tram.set({ display: display }).redraw();
          Webflow.redraw.up();
        }
      }

      // Otherwise, start a transition using the current start method.
      _tram[startMethod](clean);
      meta.start = true;
    }
  }

  // (In app) Set styles immediately and manage upstream transition
  function styleApp(el, data) {
    var _tram = tram(el);

    // Get computed transition value
    el.css('transition', '');
    var computed = el.css('transition');

    // If computed is disabled, clear upstream
    if (computed === transNone) computed = _tram.upstream = null;

    // Disable upstream temporarily
    _tram.upstream = transNone;

    // Set values immediately
    _tram.set(tramify(data));

    // Only restore upstream in preview mode
    _tram.upstream = computed;
  }

  // (Published) Set styles immediately on specified jquery element
  function stylePub(el, data) {
    tram(el).set(tramify(data));
  }

  // Build a clean object for tram
  function tramify(obj) {
    var result = {};
    var found = false;
    for (var x in obj) {
      if (x === 'transition') continue;
      result[x] = obj[x];
      found = true;
    }
    // If empty, return null for tram.set/stop compliance
    return found ? result : null;
  }

  // Events used by other webflow modules
  var events = {
    reset: function (i, el) {
      el.__wf_intro = null;
    },
    intro: function (i, el) {
      if (el.__wf_intro) return;
      el.__wf_intro = true;
      $(el).triggerHandler(introEvent);
    },
    outro: function (i, el) {
      if (!el.__wf_intro) return;
      el.__wf_intro = null;
      $(el).triggerHandler(outroEvent);
    }
  };

  // Trigger events in queue + point to sync methods
  function initEvents() {
    var count = eventQueue.length;
    for (var i = 0; i < count; i++) {
      var memo = eventQueue[i];
      memo[0](0, memo[1]);
    }
    eventQueue = [];
    $.extend(api.events, events);
  }

  // Replace events with async methods prior to init
  function asyncEvents() {
    _.each(events, function (func, name) {
      api.events[name] = function (i, el) {
        eventQueue.push([func, el]);
      };
    });
  }

  asyncEvents();

  // Export module
  return api;
});
/**
 * ----------------------------------------------------------------------
 * Webflow: Touch events
 */
Webflow.define('touch', function ($, _) {
  'use strict';

  var api = {};
  api.init = function (el) {
    el = typeof el === 'string' ? $(el).get(0) : el;
    return el ? new Touch(el) : null;
  };

  function Touch(el) {
    var active = false;
    var dirty = false;
    var useTouch = false;
    var thresholdX = Math.min(Math.round(window.innerWidth * 0.04), 40);
    var startX, startY, lastX;

    el.addEventListener('touchstart', start, false);
    el.addEventListener('touchmove', move, false);
    el.addEventListener('touchend', end, false);
    el.addEventListener('touchcancel', cancel, false);
    el.addEventListener('mousedown', start, false);
    el.addEventListener('mousemove', move, false);
    el.addEventListener('mouseup', end, false);
    el.addEventListener('mouseout', cancel, false);

    function start(evt) {
      // We don’t handle multi-touch events yet.
      var touches = evt.touches;
      if (touches && touches.length > 1) {
        return;
      }

      active = true;
      dirty = false;

      if (touches) {
        useTouch = true;
        startX = touches[0].clientX;
        startY = touches[0].clientY;
      } else {
        startX = evt.clientX;
        startY = evt.clientY;
      }

      lastX = startX;
    }

    function move(evt) {
      if (!active) return;

      if (useTouch && evt.type === 'mousemove') {
        evt.preventDefault();
        evt.stopPropagation();
        return;
      }

      var touches = evt.touches;
      var x = touches ? touches[0].clientX : evt.clientX;
      var y = touches ? touches[0].clientY : evt.clientY;

      var velocityX = x - lastX;
      lastX = x;

      if (Math.abs(velocityX) > thresholdX) {
        triggerEvent('swipe', evt, { direction: velocityX > 0 ? 'right' : 'left' });
        cancel();
      }

      // If pointer moves more than 10px flag to cancel tap
      if (Math.abs(x - startX) > 10 || Math.abs(y - startY) > 10) {
        dirty = true;
      }
    }

    function end(evt) {
      if (!active) return;

      if (useTouch && evt.type === 'mouseup') {
        evt.preventDefault();
        evt.stopPropagation();
        useTouch = false;
        return;
      }

      if (!dirty) triggerEvent('tap', evt);
    }

    function cancel(evt) {
      active = false;
    }

    function destroy() {
      el.removeEventListener('touchstart', start, false);
      el.removeEventListener('touchmove', move, false);
      el.removeEventListener('touchend', end, false);
      el.removeEventListener('touchcancel', cancel, false);
      el.removeEventListener('mousedown', start, false);
      el.removeEventListener('mousemove', move, false);
      el.removeEventListener('mouseup', end, false);
      el.removeEventListener('mouseout', cancel, false);
      el = null;
    }

    // Public instance methods
    this.destroy = destroy;
  }

  // Wrap native event to supoprt preventdefault + stopPropagation
  function triggerEvent(type, evt, data) {
    var newEvent = $.Event(type, { originalEvent: evt });
    $(evt.target).trigger(newEvent, data);
  }

  // Listen for touch events on all nodes by default.
  api.instance = api.init(document);

  // Export module
  return api;
});
/**
 * ----------------------------------------------------------------------
 * Webflow: Forms
 */
Webflow.define('forms', function ($, _) {
  'use strict';

  var api = {};
  var $doc = $(document);
  var $forms;
  var loc = window.location;
  var retro = window.XDomainRequest && !window.atob;
  var namespace = '.w-form';
  var siteId;
  var emailField = /e(\-)?mail/i;
  var emailValue = /^\S+@\S+$/;
  var alert = window.alert;

  // MailChimp domains: list-manage.com + mirrors
  var chimpRegex = /list-manage[1-9]?.com/i;

  api.ready = function () {
    // Init forms
    init();

    // Wire events
    listen && listen();
    listen = null;
  };

  api.preview = api.design = function () {
    init();
  };

  function init() {
    siteId = $('html').attr('data-wf-site');

    $forms = $(namespace + ' form');
    $forms.each(build);
  }

  function build(i, el) {
    // Store form state using namespace
    var $el = $(el);
    var data = $.data(el, namespace);
    if (!data) data = $.data(el, namespace, { form: $el }); // data.form

    reset(data);
    var wrap = $el.closest('div.w-form');
    data.done = wrap.find('> .w-form-done');
    data.fail = wrap.find('> .w-form-fail');

    var action = data.action = $el.attr('action');
    data.handler = null;
    data.redirect = $el.attr('data-redirect');

    // MailChimp form
    if (chimpRegex.test(action)) { data.handler = submitMailChimp; return; }

    // Custom form action
    if (action) return;

    // Webflow form
    if (siteId) { data.handler = submitWebflow; return; }

    // Alert for disconnected Webflow forms
    disconnected();
  }

  function listen() {
    // Handle form submission for Webflow forms
    $doc.on('submit', namespace + ' form', function(evt) {
      var data = $.data(this, namespace);
      if (data.handler) {
        data.evt = evt;
        data.handler(data);
      }
    });
  }

  // Reset data common to all submit handlers
  function reset(data) {
    var btn = data.btn = data.form.find(':input[type="submit"]');
    data.wait = data.btn.attr('data-wait') || null;
    data.success = false;
    btn.prop('disabled', false);
    data.label && btn.val(data.label);
  }

  // Disable submit button
  function disableBtn(data) {
    var btn = data.btn;
    var wait = data.wait;
    btn.prop('disabled', true);
    // Show wait text and store previous label
    if (wait) {
      data.label = btn.val();
      btn.val(wait);
    }
  }

  // Find form fields, validate, and set value pairs
  function findFields(form, result) {
    var status = null;
    result = result || {};

    // The ":input" selector is a jQuery shortcut to select all inputs, selects, textareas
    form.find(':input:not([type="submit"])').each(function(i, el) {
      var field = $(el);
      var type = field.attr('type');
      var name = field.attr('data-name') || field.attr('name') || ('Field ' + (i + 1));
      var value = field.val();

      if (type == 'checkbox') {
        value = field.is(':checked');
      } if (type == 'radio') {
        // Radio group value already processed
        if (result[name] === null || typeof result[name] == 'string') {
          return;
        }

        value = form.find('input[name="' + field.attr('name') + '"]:checked').val() || null;
      }

      if (typeof value == 'string') value = $.trim(value);
      result[name] = value;
      status = status || getStatus(field, name, value);
    });

    return status;
  }

  function getStatus(field, name, value) {
    var status = null;
    if (!field.attr('required')) return null;
    if (!value) status = 'Please fill out the required field: ' + name;
    else if (emailField.test(name) || emailField.test(field.attr('type'))) {
      if (!emailValue.test(value)) status = 'Please enter a valid email address for: ' + name;
    }
    return status;
  }

  // Submit form to Webflow
  function submitWebflow(data) {
    reset(data);

    var form = data.form;
    var payload = {
      name: form.attr('data-name') || form.attr('name') || 'Untitled Form',
      source: loc.href,
      test: Webflow.env(),
      fields: {}
    };

    preventDefault(data);

    // Find & populate all fields
    var status = findFields(form, payload.fields);
    if (status) return alert(status);

    // Disable submit button
    disableBtn(data);

    // Read site ID
    // NOTE: If this site is exported, the HTML tag must retain the data-wf-site attribute for forms to work
    if (!siteId) { afterSubmit(data); return; }
    var url = 'https://webflow.com/api/v1/form/' + siteId;

    // Work around same-protocol IE XDR limitation
    if (retro && url.indexOf('https://webflow.com') >= 0) {
      url = url.replace('https://webflow.com/', 'http://data.webflow.com/');
    }

    $.ajax({
      url: url,
      type: 'POST',
      data: payload,
      dataType: 'json',
      crossDomain: true
    }).done(function () {
      data.success = true;
      afterSubmit(data);
    }).fail(function () {
      afterSubmit(data);
    });
  }

  // Submit form to MailChimp
  function submitMailChimp(data) {
    reset(data);

    var form = data.form;
    var payload = {};

    // Skip Ajax submission if http/s mismatch, fallback to POST instead
    if (/^https/.test(loc.href) && !/^https/.test(data.action)) {
      form.attr('method', 'post');
      return;
    }

    preventDefault(data);

    // Find & populate all fields
    var status = findFields(form, payload);
    if (status) return alert(status);

    // Disable submit button
    disableBtn(data);

    // Use special format for MailChimp params
    var fullName;
    _.each(payload, function (value, key) {
      if (emailField.test(key)) payload.EMAIL = value;
      if (/^((full[ _-]?)?name)$/i.test(key)) fullName = value;
      if (/^(first[ _-]?name)$/i.test(key)) payload.FNAME = value;
      if (/^(last[ _-]?name)$/i.test(key)) payload.LNAME = value;
    });

    if (fullName && !payload.FNAME) {
      fullName = fullName.split(' ');
      payload.FNAME = fullName[0];
      payload.LNAME = payload.LNAME || fullName[1];
    }

    // Use the (undocumented) MailChimp jsonp api
    var url = data.action.replace('/post?', '/post-json?') + '&c=?';
    // Add special param to prevent bot signups
    var userId = url.indexOf('u=')+2;
    userId = url.substring(userId, url.indexOf('&', userId));
    var listId = url.indexOf('id=')+3;
    listId = url.substring(listId, url.indexOf('&', listId));
    payload['b_' + userId + '_' + listId] = '';

    $.ajax({
      url: url,
      data: payload,
      dataType: 'jsonp'
    }).done(function (resp) {
      data.success = (resp.result == 'success' || /already/.test(resp.msg));
      if (!data.success) console.info('MailChimp error: ' + resp.msg);
      afterSubmit(data);
    }).fail(function () {
      afterSubmit(data);
    });
  }

  // Common callback which runs after all Ajax submissions
  function afterSubmit(data) {
    var form = data.form;
    var wrap = form.closest('div.w-form');
    var redirect = data.redirect;
    var success = data.success;

    // Redirect to a success url if defined
    if (success && redirect) {
      Webflow.location(redirect);
      return;
    }

    // Show or hide status divs
    data.done.toggle(success);
    data.fail.toggle(!success);

    // Hide form on success
    form.toggle(!success);

    // Reset data and enable submit button
    reset(data);
  }

  function preventDefault(data) {
    data.evt && data.evt.preventDefault();
    data.evt = null;
  }

  var disconnected = _.debounce(function () {
    alert('Oops! This page has a form that is powered by Webflow, but important code was removed that is required to make the form work. Please contact support@webflow.com to fix this issue.');
  }, 100);

  // Export module
  return api;
});
/**
 * ----------------------------------------------------------------------
 * Webflow: Maps widget
 */
Webflow.define('maps', function ($, _) {
  'use strict';

  var api = {};
  var $doc = $(document);
  var google;
  var $maps;
  var namespace = '.w-widget-map';

  // -----------------------------------
  // Module methods

  api.ready = function () {
    // Init Maps on the front-end
    if (!Webflow.env()) initMaps();
  };

  api.preview = function () {
    // Update active map nodes
    $maps = $doc.find(namespace);
    // Listen for resize events
    Webflow.resize.off(triggerRedraw);
    if ($maps.length) {
      Webflow.resize.on(triggerRedraw);
      triggerRedraw();
    }
  };

  api.design = function (evt) {
    // Update active map nodes
    $maps = $doc.find(namespace);
    // Stop listening for resize events
    Webflow.resize.off(triggerRedraw);
    // Redraw to account for page changes
    $maps.length && _.defer(triggerRedraw);
  };

  // -----------------------------------
  // Private methods

  // Trigger redraw in designer or preview mode
  function triggerRedraw() {
    if ($maps.length && Webflow.app) {
      $maps.each(Webflow.app.redrawElement);
    }
  }

  function initMaps() {
    $maps = $doc.find(namespace);
    if ($maps.length) {
      Webflow.script('https://maps.googleapis.com/maps/api/js?v=3.exp&sensor=false&callback=_wf_maps_loaded');
      window._wf_maps_loaded = function () {
        window._wf_maps_loaded = function () {};
        google = window.google;
        $maps.each(renderMap);
        Webflow.resize.on(resizeMaps);
        Webflow.redraw.on(resizeMaps);
      };
    }
  }

  // Render map onto each element
  function renderMap(i, el) {
    var data = $(el).data();
    getState(el, data);
  }

  function resizeMaps() {
    $maps.each(resizeMap);
  }

  // Resize map when window changes
  function resizeMap(i, el) {
    var state = getState(el);
    google.maps.event.trigger(state.map, 'resize');
    state.setMapPosition();
  }

  // Store state on element data
  var store = 'w-widget-map';
  function getState(el, data) {

    var state = $.data(el, store);
    if (state) return state;

    var $el = $(el);
    state = $.data(el, store, {
      // Default options
      latLng: '51.511214,-0.119824',
      tooltip: '',
      style: 'roadmap',
      zoom: 12,

      // Marker
      marker: new google.maps.Marker({
        draggable: false
      }),

      // Tooltip infowindow
      infowindow: new google.maps.InfoWindow({
        disableAutoPan: true
      })
    });

    // LatLng center point
    var latLng = data.widgetLatlng || state.latLng;
    state.latLng = latLng;
    var coords = latLng.split(',');
    var latLngObj = new google.maps.LatLng(coords[0], coords[1]);
    state.latLngObj = latLngObj;

    // Disable touch events
    var mapDraggable = (Webflow.env.touch && data.disableTouch) ? false : true;

    // Map instance
    state.map = new google.maps.Map(el, {
      center: state.latLngObj,
      zoom: state.zoom,
      maxZoom: 18,
      mapTypeControl: false,
      panControl: false,
      streetViewControl: false,
      scrollwheel: !data.disableScroll,
      draggable: mapDraggable,
      zoomControl: true,
      zoomControlOptions: {
        style: google.maps.ZoomControlStyle.SMALL
      },
      mapTypeId: state.style
    });
    state.marker.setMap(state.map);

    // Set map position and offset
    state.setMapPosition = function () {
      state.map.setCenter(state.latLngObj);
      var offsetX = 0;
      var offsetY = 0;
      var padding = $el.css(['paddingTop', 'paddingRight', 'paddingBottom', 'paddingLeft']);
      offsetX -= parseInt(padding.paddingLeft, 10);
      offsetX += parseInt(padding.paddingRight, 10);
      offsetY -= parseInt(padding.paddingTop, 10);
      offsetY += parseInt(padding.paddingBottom, 10);
      if (offsetX || offsetY) {
        state.map.panBy(offsetX, offsetY);
      }
      $el.css('position', ''); // Remove injected position
    };

    // Fix position after first tiles have loaded
    google.maps.event.addListener(state.map, 'tilesloaded', function () {
      google.maps.event.clearListeners(state.map, 'tilesloaded');
      state.setMapPosition();
    });

    // Set initial position
    state.setMapPosition();
    state.marker.setPosition(state.latLngObj);
    state.infowindow.setPosition(state.latLngObj);

    // Draw tooltip
    var tooltip = data.widgetTooltip;
    if (tooltip) {
      state.tooltip = tooltip;
      state.infowindow.setContent(tooltip);
      if (!state.infowindowOpen) {
        state.infowindow.open(state.map, state.marker);
        state.infowindowOpen = true;
      }
    }

    // Map style - options.style
    var style = data.widgetStyle;
    if (style) {
      state.map.setMapTypeId(style);
    }

    // Zoom - options.zoom
    var zoom = data.widgetZoom;
    if (zoom != null) {
      state.zoom = zoom;
      state.map.setZoom(+zoom);
    }

    // Click marker to open in google maps
    google.maps.event.addListener(state.marker, 'click', function() {
      window.open('https://maps.google.com/?z=' + state.zoom + '&daddr=' + state.latLng);
    });

    return state;
  }

  // Export module
  return api;
});
/**
 * ----------------------------------------------------------------------
 * Webflow: Google+ widget
 */
Webflow.define('gplus', function ($) {
  'use strict';

  var $doc = $(document);
  var api = {};

  api.ready = function () {
    // Load Google+ API on the front-end
    if (!Webflow.env()) init();
  };

  function init() {
    $doc.find('.w-widget-gplus').length && Webflow.script('https://apis.google.com/js/plusone.js');
  }

  // Export module
  return api;
});
/**
 * ----------------------------------------------------------------------
 * Webflow: Smooth scroll
 */
Webflow.define('scroll', function ($) {
  'use strict';

  var $doc = $(document);
  var win = window;
  var loc = win.location;
  var validHash = /^[a-zA-Z][\w:.-]*$/;

  function ready() {
    // If hash is already present on page load, scroll to it right away
    if (loc.hash) {
      findEl(loc.hash.substring(1));
    }

    // When clicking on a link, check if it links to another part of the page
    $doc.on('click', 'a', function(e) {
      if (Webflow.env('design')) {
        return;
      }

      // Ignore links being used by jQuery mobile
      if (window.$.mobile && $(e.currentTarget).hasClass('ui-link')) return;

      var hash = this.hash ? this.hash.substring(1) : null;
      if (hash) {
        findEl(hash, e);
      }
    });
  }

  function findEl(hash, e) {
    if (!validHash.test(hash)) return;

    var el = $('#' + hash);
    if (!el.length) {
      return;
    }

    if (e) {
      e.preventDefault();
      e.stopPropagation();
    }

    // Push new history state
    if (loc.hash !== hash && win.history && win.history.pushState) {
      win.history.pushState(null, null, '#' + hash);
    }

    // If a fixed header exists, offset for the height
    var header = $('header, body > .header, body > .w-nav');
    var offset = header.css('position') === 'fixed' ? header.outerHeight() : 0;

    win.setTimeout(function() {
      scroll(el, offset);
    }, e ? 0 : 300);
  }

  function scroll(el, offset){
    var start = $(win).scrollTop();
    var end = el.offset().top - offset;

    // If specified, scroll so that the element ends up in the middle of the viewport
    if (el.data('scroll') == 'mid') {
      var available = $(win).height() - offset;
      var elHeight = el.outerHeight();
      if (elHeight < available) {
        end -= Math.round((available - elHeight) / 2);
      }
    }

    var mult = 1;

    // Check for custom time multiplier on the body and the element
    $('body').add(el).each(function(i) {
      var time = parseFloat($(this).attr('data-scroll-time'), 10);
      if (!isNaN(time) && (time === 0 || time > 0)) {
        mult = time;
      }
    });

    // Shim for IE8 and below
    if (!Date.now) {
      Date.now = function() { return new Date().getTime(); };
    }

    var clock = Date.now();
    var animate = win.requestAnimationFrame || win.mozRequestAnimationFrame || win.webkitRequestAnimationFrame || function(fn) { win.setTimeout(fn, 15); };
    var duration = (472.143 * Math.log(Math.abs(start - end) +125) - 2000) * mult;

    var step = function() {
      var elapsed = Date.now() - clock;
      win.scroll(0, getY(start, end, elapsed, duration));

      if (elapsed <= duration) {
        animate(step);
      }
    };

    step();
  }

  function getY(start, end, elapsed, duration) {
    if (elapsed > duration) {
      return end;
    }

    return start + (end - start) * ease(elapsed / duration);
  }

  function ease(t) {
    return t<0.5 ? 4*t*t*t : (t-1)*(2*t-2)*(2*t-2)+1;
  }

  // Export module
  return { ready: ready };
});
/**
 * ----------------------------------------------------------------------
 * Webflow: Auto-select links to current page or section
 */
Webflow.define('links', function ($, _) {
  'use strict';

  var api = {};
  var $win = $(window);
  var designer;
  var inApp = Webflow.env();
  var location = window.location;
  var linkCurrent = 'w--current';
  var validHash = /^#[a-zA-Z][\w:.-]*$/;
  var indexPage = /index\.(html|php)$/;
  var dirList = /\/$/;
  var anchors;

  // -----------------------------------
  // Module methods

  api.ready = api.design = api.preview = init;

  // -----------------------------------
  // Private methods

  function init() {
    designer = inApp && Webflow.env('design');

    // Reset scroll listener, init anchors
    Webflow.scroll.off(scroll);
    anchors = [];

    // Test all links for a selectable href
    var links = document.links;
    for (var i = 0; i < links.length; ++i) {
      select(links[i]);
    }

    // Listen for scroll if any anchors exist
    if (anchors.length) {
      Webflow.scroll.on(scroll);
      scroll();
    }
  }

  function select(link) {
    var href = link.getAttribute('href');

    // Ignore any hrefs with a colon to safely avoid all uri schemes
    if (href.indexOf(':') >= 0) return;

    var $link = $(link);

    // Check for valid hash links w/ sections and use scroll anchor
    if (href.indexOf('#') === 0 && validHash.test(href)) {
      var $section = $(href);
      $section.length && anchors.push({ link: $link, sec: $section, active: false });
      return;
    }

    // Determine whether the link should be selected
    var slug = (inApp ? Webflow.env('slug') : location.pathname) || '';
    var match = (link.href === location.href) || (href === slug) || (indexPage.test(href) && dirList.test(slug));
    setClass($link, linkCurrent, match);
  }

  function scroll() {
    var viewTop = $win.scrollTop();
    var viewHeight = $win.height();

    // Check each anchor for a section in view
    _.each(anchors, function (anchor) {
      var $link = anchor.link;
      var $section = anchor.sec;
      var top = $section.offset().top;
      var height = $section.outerHeight();
      var offset = viewHeight * 0.5;
      var active = ($section.is(':visible') &&
        top + height - offset >= viewTop &&
        top + offset <= viewTop + viewHeight);
      if (anchor.active === active) return;
      anchor.active = active;
      setClass($link, linkCurrent, active);
      if (designer) $link[0].__wf_current = active;
    });
  }

  function setClass($elem, className, add) {
    var exists = $elem.hasClass(className);
    if (add && exists) return;
    if (!add && !exists) return;
    add ? $elem.addClass(className) : $elem.removeClass(className);
  }

  // Export module
  return api;
});
/**
 * ----------------------------------------------------------------------
 * Webflow: Slider component
 */
Webflow.define('slider', function ($, _) {
  'use strict';

  var api = {};
  var tram = window.tram;
  var $doc = $(document);
  var $sliders;
  var designer;
  var inApp = Webflow.env();
  var namespace = '.w-slider';
  var dot = '<div class="w-slider-dot" data-wf-ignore />';
  var ix = Webflow.require('ix');
  ix = ix && ix.events;
  var fallback;
  var redraw;

  // -----------------------------------
  // Module methods

  api.ready = function () {
    init();
  };

  api.design = function () {
    designer = true;
    init();
  };

  api.preview = function () {
    designer = false;
    init();
  };

  api.redraw = function () {
    redraw = true;
    init();
  };

  // -----------------------------------
  // Private methods

  function init() {
    // Find all sliders on the page
    $sliders = $doc.find(namespace);
    $sliders.each(build);
    redraw = null;
    if (fallback) return;

    // Wire events
    listen && listen();
    listen = null;
  }

  function listen() {
    Webflow.resize.on(function () {
      $sliders.each(render);
    });

    Webflow.redraw.on(api.redraw);
  }

  function build(i, el) {
    var $el = $(el);

    // Store slider state in data
    var data = $.data(el, namespace);
    if (!data) data = $.data(el, namespace, { index: 0, el: $el, config: {} });
    data.mask = $el.children('.w-slider-mask');
    data.left = $el.children('.w-slider-arrow-left');
    data.right = $el.children('.w-slider-arrow-right');
    data.nav = $el.children('.w-slider-nav');
    data.slides = data.mask.children('.w-slide');
    if (ix) data.slides.each(ix.reset);
    if (redraw) data.maskWidth = 0;

    // Disable in old browsers
    if (!tram.support.transform) {
      data.left.hide();
      data.right.hide();
      data.nav.hide();
      fallback = true;
      return;
    }

    // Remove old events
    data.el.off(namespace);
    data.left.off(namespace);
    data.right.off(namespace);
    data.nav.off(namespace);

    // Set config from data attributes
    configure(data);

    // Add events based on mode
    if (designer) {
      data.el.on('setting' + namespace, handler(data));
      killTimer(data);
      data.hasTimer = false;
    } else {
      data.el.on('swipe' + namespace, handler(data));
      data.left.on('tap' + namespace, previous(data));
      data.right.on('tap' + namespace, next(data));

      // Start timer if autoplay is true, only once
      if (data.config.autoplay && !data.hasTimer) {
        data.hasTimer = true;
        startTimer(data);
      }
    }

    // Listen to nav events
    data.nav.on('tap' + namespace, '> div', handler(data));

    // Remove gaps from formatted html (for inline-blocks)
    if (!inApp) {
      data.mask.contents().filter(function() {
        return this.nodeType === 3;
      }).remove();
    }

    // Run first render
    render(i, el);
  }

  function configure(data) {
    var config = {};

    config.depth = 1;
    config.crossOver = 0;

    // Set config options from data attributes
    config.animation = data.el.attr('data-animation') || 'slide';
    if (config.animation == 'outin') {
      config.animation = 'cross';
      config.crossOver = 0.5;
    }
    config.easing = data.el.attr('data-easing') || 'ease';

    var duration = data.el.attr('data-duration');
    config.duration = duration != null ? +duration : 500;

    if (+data.el.attr('data-infinite')) config.infinite = true;

    if (+data.el.attr('data-hide-arrows')) {
      config.hideArrows = true;
    } else if (data.config.hideArrows) {
      data.left.show();
      data.right.show();
    }

    if (+data.el.attr('data-autoplay')) {
      config.autoplay = true;
      config.delay = +data.el.attr('data-delay') || 2000;
      // Disable timer on first touch or mouse down
      var touchEvents = 'mousedown' + namespace + ' touchstart' + namespace;
      if (!designer) data.el.off(touchEvents).one(touchEvents, function () {
        killTimer(data);
      });
    }

    // Use edge buffer to help calculate page count
    var arrowWidth = data.right.width();
    config.edge = arrowWidth ? arrowWidth + 40 : 100;

    // Store config in data
    data.config = config;
  }

  function previous(data) {
    return function (evt) {
      change(data, { index: data.index - 1, vector: -1 });
    };
  }

  function next(data) {
    return function (evt) {
      change(data, { index: data.index + 1, vector: 1 });
    };
  }

  function select(data, value) {
    // Select page based on slide element index
    var found = null;
    if (value === data.slides.length) {
      init(); layout(data); // Rebuild and find new slides
    }
    _.each(data.anchors, function (anchor, index) {
      $(anchor.els).each(function (i, el) {
        if ($(el).index() === value) found = index;
      });
    });
    if (found != null) change(data, { index: found, immediate: true });
  }

  function startTimer(data) {
    var config = data.config;
    stopTimer(data);
    config.timer = window.setTimeout(function () {
      if (!config.autoplay || designer) return;
      next(data)();
      startTimer(data);
    }, config.delay);
  }

  function stopTimer(data) {
    var config = data.config;
    window.clearTimeout(config.timer);
    config.timer = null;
  }

  function killTimer(data) {
    var config = data.config;
    config.autoplay = false;
    stopTimer(data);
  }

  function handler(data) {
    return function (evt, options) {
      options = options || {};

      // Designer settings
      if (designer && evt.type == 'setting') {
        if (options.select == 'prev') return previous(data)();
        if (options.select == 'next') return next(data)();
        configure(data);
        layout(data);
        if (options.select == null) return;
        select(data, options.select);
        return;
      }

      // Swipe event
      if (evt.type == 'swipe') {
        if (options.direction == 'left') return next(data)();
        if (options.direction == 'right') return previous(data)();
        return;
      }

      // Page buttons
      if (data.nav.has(evt.target).length) {
        change(data, { index: $(evt.target).index() });
      }
    };
  }

  function change(data, options) {
    options = options || {};
    var config = data.config;
    var anchors = data.anchors;

    // Set new index
    data.previous = data.index;
    var index = options.index;
    var shift = {};
    if (index < 0) {
      index = anchors.length-1;
      if (config.infinite) {
        // Shift first slide to the end
        shift.x = -data.endX;
        shift.from = 0;
        shift.to = anchors[0].width;
      }
    } else if (index >= anchors.length) {
      index = 0;
      if (config.infinite) {
        // Shift last slide to the start
        shift.x = anchors[anchors.length-1].width;
        shift.from = -anchors[anchors.length-1].x;
        shift.to = shift.from - shift.x;
      }
    }
    data.index = index;

    // Select page nav
    var active = data.nav.children().eq(data.index).addClass('w-active');
    data.nav.children().not(active).removeClass('w-active');

    // Hide arrows
    if (config.hideArrows) {
      data.index === anchors.length-1 ? data.right.hide() : data.right.show();
      data.index === 0 ? data.left.hide() : data.left.show();
    }

    // Get page offset from anchors
    var lastOffsetX = data.offsetX || 0;
    var offsetX = data.offsetX = -anchors[data.index].x;
    var resetConfig = { x: offsetX, opacity: 1, visibility: '' };

    // Transition slides
    var targets = $(anchors[data.index].els);
    var previous = $(anchors[data.previous] && anchors[data.previous].els);
    var others = data.slides.not(targets);
    var animation = config.animation;
    var easing = config.easing;
    var duration = Math.round(config.duration);
    var vector = options.vector || (data.index > data.previous ? 1 : -1);
    var fadeRule = 'opacity ' + duration + 'ms ' + easing;
    var slideRule = 'transform ' + duration + 'ms ' + easing;

    // Trigger IX events
    if (!designer && ix) {
      targets.each(ix.intro);
      others.each(ix.outro);
    }

    // Set immediately after layout changes (but not during redraw)
    if (options.immediate && !redraw) {
      tram(targets).set(resetConfig);
      resetOthers();
      return;
    }

    // Exit early if index is unchanged
    if (data.index == data.previous) return;

    // Cross Fade / Out-In
    if (animation == 'cross') {
      var reduced = Math.round(duration - duration * config.crossOver);
      var wait = Math.round(duration - reduced);
      fadeRule = 'opacity ' + reduced + 'ms ' + easing;
      tram(previous)
        .set({ visibility: '' })
        .add(fadeRule)
        .start({ opacity: 0 });
      tram(targets)
        .set({ visibility: '', x: offsetX, opacity: 0, zIndex: config.depth++ })
        .add(fadeRule)
        .wait(wait)
        .then({ opacity: 1 })
        .then(resetOthers);
      return;
    }

    // Fade Over
    if (animation == 'fade') {
      tram(previous)
        .set({ visibility: '' })
        .stop();
      tram(targets)
        .set({ visibility: '', x: offsetX, opacity: 0, zIndex: config.depth++ })
        .add(fadeRule)
        .start({ opacity: 1 })
        .then(resetOthers);
      return;
    }

    // Slide Over
    if (animation == 'over') {
      resetConfig = { x: data.endX };
      tram(previous)
        .set({ visibility: '' })
        .stop();
      tram(targets)
        .set({ visibility: '', zIndex: config.depth++, x: offsetX + anchors[data.index].width * vector })
        .add(slideRule)
        .start({ x: offsetX })
        .then(resetOthers);
      return;
    }

    // Slide - infinite scroll
    if (config.infinite && shift.x) {
      tram(data.slides.not(previous))
        .set({ visibility: '', x: shift.x })
        .add(slideRule)
        .start({ x: offsetX });
      tram(previous)
        .set({ visibility: '', x: shift.from })
        .add(slideRule)
        .start({ x: shift.to });
      data.shifted = previous;

    } else {
      if (config.infinite && data.shifted) {
        tram(data.shifted).set({ visibility: '', x: lastOffsetX });
        data.shifted = null;
      }

      // Slide - basic scroll
      tram(data.slides)
        .set({ visibility: '' })
        .add(slideRule)
        .start({ x: offsetX });
    }

    // Helper to move others out of view
    function resetOthers() {
      var targets = $(anchors[data.index].els);
      var others = data.slides.not(targets);
      if (animation != 'slide') resetConfig.visibility = 'hidden';
      tram(others).set(resetConfig);
    }
  }

  function render(i, el) {
    var data = $.data(el, namespace);
    if (maskChanged(data)) return layout(data);
    if (designer && slidesChanged(data)) layout(data);
  }

  function layout(data) {
    // Determine page count from width of slides
    var pages = 1;
    var offset = 0;
    var anchor = 0;
    var width = 0;
    data.anchors = [{ els: [], x: 0, width: 0 }];
    data.slides.each(function (i, el) {
      if (anchor - offset > data.maskWidth - data.config.edge) {
        pages++;
        offset += data.maskWidth;
        // Store page anchor for transition
        data.anchors[pages-1] = { els: [], x: anchor, width: 0 };
      }
      // Set next anchor using current width + margin
      width = $(el).outerWidth(true);
      anchor += width;
      data.anchors[pages-1].width += width;
      data.anchors[pages-1].els.push(el);
    });
    data.endX = anchor;

    // Build dots if nav exists and needs updating
    if (designer) data.pages = null;
    if (data.nav.length && data.pages !== pages){
      data.pages = pages;
      buildNav(data);
    }

    // Make sure index is still within range and call change handler
    var index = data.index;
    if (index >= pages) index = pages-1;
    change(data, { immediate: true, index: index });
  }

  function buildNav(data) {
    var dots = [];
    var $dot;
    var spacing = data.el.attr('data-nav-spacing');
    if (spacing) spacing = parseFloat(spacing) + 'px';
    for (var i=0; i<data.pages; i++) {
      $dot = $(dot);
      if (data.nav.hasClass('w-num')) $dot.text(i+1);
      if (spacing != null) $dot.css({
        'margin-left': spacing,
        'margin-right': spacing
      });
      dots.push($dot);
    }
    data.nav.empty().append(dots);
  }

  function maskChanged(data) {
    var maskWidth = data.mask.width();
    if (data.maskWidth !== maskWidth) {
      data.maskWidth = maskWidth;
      return true;
    }
    return false;
  }

  function slidesChanged(data) {
    var slidesWidth = 0;
    data.slides.each(function (i, el) {
      slidesWidth += $(el).outerWidth(true);
    });
    if (data.slidesWidth !== slidesWidth) {
      data.slidesWidth = slidesWidth;
      return true;
    }
    return false;
  }

  // Export module
  return api;
});
/**
 * ----------------------------------------------------------------------
 * Webflow: Navbar component
 */
Webflow.define('navbar', function ($, _) {
  'use strict';

  var api = {};
  var tram = window.tram;
  var $win = $(window);
  var $doc = $(document);
  var $body;
  var $navbars;
  var designer;
  var inApp = Webflow.env();
  var overlay = '<div class="w-nav-overlay" data-wf-ignore />';
  var namespace = '.w-nav';
  var buttonOpen = 'w--open';
  var menuOpen = 'w--nav-menu-open';
  var linkOpen = 'w--nav-link-open';

  // -----------------------------------
  // Module methods

  api.ready = api.design = api.preview = init;

  // -----------------------------------
  // Private methods

  function init() {
    designer = inApp && Webflow.env('design');
    $body = $(document.body);

    // Find all instances on the page
    $navbars = $doc.find(namespace);
    $navbars.each(build);

    // Wire events
    listen && listen();
    listen = null;
  }

  function listen() {
    Webflow.resize.on(function () {
      $navbars.each(resize);
    });
  }

  function build(i, el) {
    var $el = $(el);

    // Store state in data
    var data = $.data(el, namespace);
    if (!data) data = $.data(el, namespace, { open: false, el: $el, config: {} });
    data.menu = $el.find('.w-nav-menu');
    data.links = data.menu.find('.w-nav-link');
    data.button = $el.find('.w-nav-button');
    data.container = $el.find('.w-container');
    data.outside = outside(data);

    // Remove old events
    data.el.off(namespace);
    data.button.off(namespace);
    data.menu.off(namespace);

    // Set config from data attributes
    configure(data);

    // Add events based on mode
    if (designer) {
      removeOverlay(data);
      data.el.on('setting' + namespace, handler(data));
    } else {
      addOverlay(data);
      data.button.on('tap' + namespace, toggle(data));
      data.menu.on('tap' + namespace, 'a', navigate(data));
    }

    // Trigger initial resize
    resize(i, el);
  }

  function removeOverlay(data) {
    if (!data.overlay) return;
    close(data, true);
    data.overlay.remove();
    data.overlay = null;
  }

  function addOverlay(data) {
    if (data.overlay) return;
    data.overlay = $(overlay).appendTo(data.el);
    data.parent = data.menu.parent();
    close(data, true);
  }

  function configure(data) {
    var config = {};
    var old = data.config || {};

    // Set config options from data attributes
    config.animation = data.el.attr('data-animation') || 'default';

    // Re-open menu if the animation type changed
    if (old.animation != config.animation) {
      data.open && _.defer(reopen, data);
    }

    config.easing = data.el.attr('data-easing') || 'ease';
    config.easing2 = data.el.attr('data-easing2') || 'ease';

    var duration = data.el.attr('data-duration');
    config.duration = duration != null ? +duration : 400;

    config.docHeight = data.el.attr('data-doc-height');

    // Store config in data
    data.config = config;
  }

  function handler(data) {
    return function (evt, options) {
      options = options || {};

      // Designer settings
      if (designer && evt.type == 'setting') {
        var winWidth = $win.width();
        configure(data);
        options.open === true && open(data, true);
        options.open === false && close(data, true);
        // Reopen if media query changed after setting
        data.open && _.defer(function () {
          if (winWidth != $win.width()) reopen(data);
        });
        return;
      }
    };
  }

  function closeEach(i, el) {
    var data = $.data(el, namespace);
    data.open && close(data);
  }

  function reopen(data) {
    if (!data.open) return;
    close(data, true);
    open(data, true);
  }

  function toggle(data) {
    return _.debounce(function (evt) {
      data.open ? close(data) : open(data);
    });
  }

  function navigate(data) {
    return function (evt) {
      var link = $(this);
      var href = link.attr('href');

      // Close when navigating to an in-page anchor
      if (href && href.indexOf('#') === 0 && data.open) {
        // Avoid empty hash links
        if (href.length === 1) evt.preventDefault();
        close(data);
      }
    };
  }

  function outside(data) {
    // Close menu when tapped outside
    return _.debounce(function (evt) {
      if (!data.open) return;
      var menu = $(evt.target).closest('.w-nav-menu');
      if (!data.menu.is(menu)) {
        close(data);
      }
    });
  }

  function resize(i, el) {
    var data = $.data(el, namespace);
    // Check for collapsed state based on button display
    var collapsed = data.collapsed = data.button.css('display') != 'none';
    // Close menu if button is no longer visible (and not in designer)
    if (data.open && !collapsed && !designer) close(data, true);
    // Set max-width of links to match container
    data.container.length && data.links.each(maxLink(data));
    // If currently open and in overlay mode, update height to match body
    if (data.open && /^over/.test(data.config.animation)) {
      setBodyHeight(data);
      setMenuHeight(data);
    }
  }

  var maxWidth = 'max-width';
  function maxLink(data) {
    // Set max-width of each link (unless it has an upstream value)
    var containMax = data.container.css(maxWidth);
    if (containMax == 'none') containMax = '';
    return function (i, link) {
      link = $(link);
      link.css(maxWidth, '');
      if (link.css(maxWidth) == 'none') link.css(maxWidth, containMax);
    };
  }

  function open(data, immediate) {
    if (data.open) return;
    data.open = true;
    data.menu.addClass(menuOpen);
    data.links.addClass(linkOpen);
    data.button.addClass(buttonOpen);
    var config = data.config;
    var animation = config.animation;
    if (animation == 'none' || !tram.support.transform) immediate = true;
    var animOver = /^over/.test(animation);
    var bodyHeight = setBodyHeight(data);
    var menuHeight = data.menu.outerHeight(true);
    var menuWidth = data.menu.outerWidth(true);
    var navHeight = data.el.height();
    var direction = /left$/.test(animation) ? -1 : 1;
    resize(0, data.el[0]);

    // Listen for tap outside events
    if (!designer) $doc.on('tap' + namespace, data.outside);

    // Update menu height for Over state
    if (animOver) setMenuHeight(data);

    // No transition for immediate
    if (immediate) return;

    var transConfig = 'transform ' + config.duration + 'ms ' + config.easing;

    // Add menu to overlay
    if (data.overlay) {
      data.overlay.show()
        .append(data.menu)
        .height(menuHeight);
    }

    // Over left/right
    if (animOver) {
      tram(data.menu)
        .add(transConfig)
        .set({ x: direction * menuWidth, height: bodyHeight }).start({ x: 0 });
      data.overlay && data.overlay.css({ width: menuWidth, height: bodyHeight });
      return;
    }

    // Drop Down
    var offsetY = navHeight + menuHeight;
    tram(data.menu)
      .add(transConfig)
      .set({ y: -offsetY }).start({ y: 0 });
  }

  function setBodyHeight(data) {
    return data.bodyHeight = data.config.docHeight ? $doc.height() : $body.height();
  }

  function setMenuHeight(data) {
    var bodyHeight = data.bodyHeight;
    var navFixed = data.el.css('position') == 'fixed';
    if (!navFixed) bodyHeight -= data.el.offset().top;
    data.menu.height(bodyHeight);
  }

  function close(data, immediate) {
    data.open = false;
    data.button.removeClass(buttonOpen);
    var config = data.config;
    if (config.animation == 'none' || !tram.support.transform) immediate = true;
    var animation = config.animation;

    // Stop listening for tap outside events
    $doc.off('tap' + namespace, data.outside);

    if (immediate) {
      tram(data.menu).stop();
      complete();
      return;
    }

    var transConfig = 'transform ' + config.duration + 'ms ' + config.easing2;
    var menuHeight = data.menu.outerHeight(true);
    var menuWidth = data.menu.outerWidth(true);
    var navHeight = data.el.height();
    var direction = /left$/.test(animation) ? -1 : 1;
    var animOver = /^over/.test(animation);

    // Over left/right
    if (animOver) {
      tram(data.menu)
        .add(transConfig)
        .start({ x: menuWidth * direction }).then(complete);
      return;
    }

    // Drop Down
    var offsetY = navHeight + menuHeight;
    tram(data.menu)
      .add(transConfig)
      .start({ y: -offsetY }).then(complete);

    function complete() {
      data.menu.height('');
      tram(data.menu).set({ x: 0, y: 0 });
      data.menu.removeClass(menuOpen);
      data.links.removeClass(linkOpen);
      if (data.overlay && data.overlay.children().length) {
        // Move menu back to parent
        data.menu.appendTo(data.parent);
        data.overlay.attr('style', '').hide();
      }
    }
  }

  // Export module
  return api;
});
/**
 * ----------------------------------------------------------------------
 * Webflow: Tabs component
 */
Webflow.define('tabs', function ($, _) {
  'use strict';

  var api = {};
  var tram = window.tram;
  var $win = $(window);
  var $doc = $(document);
  var $tabs;
  var design;
  var env = Webflow.env;
  var safari = env.safari;
  var inApp = env();
  var tabAttr = 'data-w-tab';
  var namespace = '.w-tabs';
  var linkCurrent = 'w--current';
  var tabActive = 'w--tab-active';
  var ix = Webflow.require('ix').events;

  // -----------------------------------
  // Module methods

  api.ready = api.design = api.preview = init;

  // -----------------------------------
  // Private methods

  function init() {
    design = inApp && Webflow.env('design');

    // Find all instances on the page
    $tabs = $doc.find(namespace);
    $tabs.each(build);
  }

  function build(i, el) {
    var $el = $(el);

    // Store state in data
    var data = $.data(el, namespace);
    if (!data) data = $.data(el, namespace, { el: $el, config: {} });
    data.current = null;
    data.menu = $el.children('.w-tab-menu');
    data.links = data.menu.children('.w-tab-link');
    data.content = $el.children('.w-tab-content');
    data.panes = data.content.children('.w-tab-pane');

    // Remove old events
    data.el.off(namespace);
    data.links.off(namespace);

    // Set config from data attributes
    configure(data);

    // Wire up events when not in design mode
    if (!design) {
      data.links.on('click' + namespace, linkSelect(data));

      // Trigger first intro event from current tab
      var $link = data.links.filter('.' + linkCurrent);
      var tab = $link.attr(tabAttr);
      tab && changeTab(data, { tab: tab, immediate: true });
    }
  }

  function configure(data) {
    var config = {};
    var old = data.config || {};

    // Set config options from data attributes
    config.easing = data.el.attr('data-easing') || 'ease';

    var intro = +data.el.attr('data-duration-in');
    intro = config.intro = intro === intro ? intro : 0;

    var outro = +data.el.attr('data-duration-out');
    outro = config.outro = outro === outro ? outro : 0;

    config.immediate = !intro && !outro;

    // Store config in data
    data.config = config;
  }

  function linkSelect(data) {
    return function (evt) {
      var tab = evt.currentTarget.getAttribute(tabAttr);
      tab && changeTab(data, { tab: tab });
    };
  }

  function changeTab(data, options) {
    options = options || {};

    var config = data.config;
    var easing = config.easing;
    var tab = options.tab;

    // Don't select the same tab twice
    if (tab === data.current) return;
    data.current = tab;

    // Select the current link
    data.links.each(function (i, el) {
      var $el = $(el);
      if (el.getAttribute(tabAttr) === tab) $el.addClass(linkCurrent).each(ix.intro);
      else if ($el.hasClass(linkCurrent)) $el.removeClass(linkCurrent).each(ix.outro);
    });

    // Find the new tab panes and keep track of previous
    var targets = [];
    var previous = [];
    data.panes.each(function (i, el) {
      var $el = $(el);
      if (el.getAttribute(tabAttr) === tab) {
        targets.push(el);
      } else if ($el.hasClass(tabActive)) {
        previous.push(el);
      }
    });

    var $targets = $(targets);
    var $previous = $(previous);

    // Switch tabs immediately and bypass transitions
    if (options.immediate || config.immediate) {
      $targets.addClass(tabActive).each(ix.intro);
      $previous.removeClass(tabActive);
      Webflow.redraw.up();
      return;
    }

    // Fade out the currently active tab before intro
    if ($previous.length && config.outro) {
      $previous.each(ix.outro);
      tram($previous)
        .add('opacity ' + config.outro + 'ms ' + easing, { fallback: safari })
        .start({ opacity: 0 })
        .then(intro);
    } else {
      // Skip the outro and play intro
      intro();
    }

    // Fade in the new target
    function intro() {
      // Clear previous active class + inline style
      $previous.removeClass(tabActive).removeAttr('style');

      // Add active class to new target
      $targets.addClass(tabActive).each(ix.intro);
      Webflow.redraw.up();

      // Set opacity immediately if intro is zero
      if (!config.intro) return tram($targets).set({ opacity: 1 });

      // Otherwise fade in opacity
      tram($targets)
        .set({ opacity: 0 })
        .redraw()
        .add('opacity ' + config.intro + 'ms ' + easing, { fallback: safari })
        .start({ opacity: 1 });
    }
  }

  // Export module
  return api;
});
/**
 * ----------------------------------------------------------------------
 * Webflow: Interactions: Init
 */
Webflow.require('ix').init([
  {"slug":"slicefromleft","name":"SliceFromLeft","value":{"style":{"opacity":0,"x":"-100px","y":"0px"},"triggers":[{"type":"scroll","stepsA":[{"opacity":1,"transition":"transform 500ms ease 0ms, opacity 500ms ease 0ms","x":"0px","y":"0px"}],"stepsB":[]}]}},
  {"slug":"fadein","name":"FadeIn","value":{"style":{"opacity":0,"scale":0.5,"rotate":"-25deg"},"triggers":[{"type":"load","stepsA":[{"opacity":1,"transition":"transform 500ms ease-out-cubic 0ms, opacity 500ms ease 0ms","scale":1.1,"rotate":"10deg"},{"transition":"transform 300ms ease-in-out-cubic 0ms","scale":1,"rotate":"0deg"}],"stepsB":[]}]}},
  {"slug":"shaking","name":"Shaking","value":{"style":{},"triggers":[{"type":"hover","loopA":true,"stepsA":[{"transition":"transform 200ms ease-in-out-sine 0ms","rotate":"-5deg"},{"transition":"transform 200ms ease-in-out-sine 0ms","rotate":"5deg"},{"transition":"transform 200ms ease-in-out-sine 0ms","rotate":"-5deg"},{"wait":4000,"transition":"transform 125ms ease-in-out-sine 0ms","rotate":"0deg"}],"stepsB":[{"transition":"transform 500ms ease 0ms","rotate":"0deg","scale":1}]}]}},
  {"slug":"slicefromright","name":"SliceFromRight","value":{"style":{"opacity":0,"x":"100px","y":"0px"},"triggers":[{"type":"scroll","stepsA":[{"opacity":1,"transition":"transform 500ms ease 0ms, opacity 500ms ease 0ms","x":"0px","y":"0px"}],"stepsB":[]}]}}
]);
