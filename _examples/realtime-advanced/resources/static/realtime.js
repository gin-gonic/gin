

function StartRealtime(roomid, timestamp) {
    StartEpoch(timestamp);
    StartSSE(roomid);
    StartForm();
}

function StartForm() {
    $('#chat-message').focus();
    $('#chat-form').ajaxForm(function() {
        $('#chat-message').val('');
        $('#chat-message').focus();
    });
}

function StartEpoch(timestamp) {
    var windowSize = 60;
    var height = 200;
    var defaultData = histogram(windowSize, timestamp);

    window.heapChart = $('#heapChart').epoch({
        type: 'time.area',
        axes: ['bottom', 'left'],
        height: height,
        historySize: 10,
        data: [
            {values: defaultData},
            {values: defaultData}
        ]
    });

    window.mallocsChart = $('#mallocsChart').epoch({
        type: 'time.area',
        axes: ['bottom', 'left'],
        height: height,
        historySize: 10,
        data: [
            {values: defaultData},
            {values: defaultData}
        ]
    });

    window.messagesChart = $('#messagesChart').epoch({
        type: 'time.line',
        axes: ['bottom', 'left'],
        height: 240,
        historySize: 10,
        data: [
            {values: defaultData},
            {values: defaultData},
            {values: defaultData}
        ]
    });
}

function StartSSE(roomid) {
    if (!window.EventSource) {
        alert("EventSource is not enabled in this browser");
        return;
    }
    var source = new EventSource('/stream/'+roomid);
    source.addEventListener('message', newChatMessage, false);
    source.addEventListener('stats', stats, false);
}

function stats(e) {
    var data = parseJSONStats(e.data);
    heapChart.push(data.heap);
    mallocsChart.push(data.mallocs);
    messagesChart.push(data.messages);
}

function parseJSONStats(e) {
    var data = jQuery.parseJSON(e);
    var timestamp = data.timestamp;

    var heap = [
        {time: timestamp, y: data.HeapInuse},
        {time: timestamp, y: data.StackInuse}
    ];

    var mallocs = [
        {time: timestamp, y: data.Mallocs},
        {time: timestamp, y: data.Frees}
    ];
    var messages = [
        {time: timestamp, y: data.Connected},
        {time: timestamp, y: data.Inbound},
        {time: timestamp, y: data.Outbound}
    ];

    return {
        heap: heap,
        mallocs: mallocs,
        messages: messages
    }
}

function newChatMessage(e) {
    var data = jQuery.parseJSON(e.data);
    var nick = data.nick;
    var message = data.message;
    var style = rowStyle(nick);
    var html = "<tr class=\""+style+"\"><td>"+nick+"</td><td>"+message+"</td></tr>";
    $('#chat').append(html);

    $("#chat-scroll").scrollTop($("#chat-scroll")[0].scrollHeight);
}

function histogram(windowSize, timestamp) {
    var entries = new Array(windowSize);
    for(var i = 0; i < windowSize; i++) {
        entries[i] = {time: (timestamp-windowSize+i-1), y:0};
    }
    return entries;
}

var entityMap = {
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
    '"': '&quot;',
    "'": '&#39;',
    "/": '&#x2F;'
};

function rowStyle(nick) {
    var classes = ['active', 'success', 'info', 'warning', 'danger'];
    var index = hashCode(nick)%5;
    return classes[index];
}

function hashCode(s){
  return Math.abs(s.split("").reduce(function(a,b){a=((a<<5)-a)+b.charCodeAt(0);return a&a},0));             
}

function escapeHtml(string) {
    return String(string).replace(/[&<>"'\/]/g, function (s) {
      return entityMap[s];
    });
}

window.StartRealtime = StartRealtime
