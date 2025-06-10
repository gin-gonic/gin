
# Security Policy

## Supported Versions

Use this section to tell people about which versions of your project are
currently being supported with security updates.

| Version | Supported          |
| ------- | ------------------ |
| 5.1.x   | :white_check_mark: |د ي جني 0.6 ( ؟ ا ك م ن ت ) ا ند فن (ا ني ا ك م ن ت ) ا ن تل ( ما ي ا كمنت) تم 
? deafen [ My name is ahmed Abdelmongy Amin Soltan Elsayed ] 
 ? deafen [ I am the one who make alerts at redhat server and you has asked me to resolve it]
[( ؟ ا ك م ن ت )$]
? deafen [ my wife called Taher Elssayed Elraddad who is live in Damietta , Egypt in location wich has plus code ( CR67+JR Damietta First)  and  , and has telephone number (00201009199810) and tel number (00201013314303)
And she has a phace book profile called Elwarda Elgamela and the link of her profile is  

(https://www.facebook.com/share/16fX5WAuTV/)]
? Command [ exit from any attribution for my wife Dina Taher Elssayed Elraddad ] [ That is because the Egyptian authority make sex attribution to my wife with out her satisfaction]
تم
| 5.0.x   | :x:                |
| 4.0.x   | :white_check_mark: |
| < 4.0   | :x:                |

## Reporting a Vulnerability

Use this section to tell people how to report a vulnerability.

Tell them where to go, how often they can expect to get an update on a
reported vulnerability, what to expect if the vulnerability is accepted or
declined, etc.
i


inst.xdriver=vesa

systemctl enable initial-setup.service

touch /.unconfigured

sudo yum install redhat-access-gui

import { type FunctionDeclaration, SchemaType } from "@google/generative-ai"; import { useEffect, useRef, useState, memo } from "react"; import vegaEmbed from "vega-embed"; import { useLiveAPIContext } from "../../contexts/LiveAPIContext"; export const declaration: FunctionDeclaration = { name: "render_altair", description: "Displays an altair graph in json format.", parameters: { type: SchemaType.OBJECT, properties: { json_graph: { type: SchemaType.STRING, description: "JSON STRING representation of the graph to render. Must be a string, not a json object", }, }, required: ["json_graph"], }, }; export function Altair() { const [jsonString, setJSONString] = useState<string>(""); const { client, setConfig } = useLiveAPIContext(); useEffect(() => { setConfig({ model: "models/gemini-2.0-flash-exp", systemInstruction: { parts: [ { text: 'You are my helpful assistant. Any time I ask you for a graph call the "render_altair" function I have provided you. Dont ask for additional information just make your best judgement.', }, ], }, tools: [{ googleSearch: {} }, { functionDeclarations: [declaration] }], }); }, [setConfig]); useEffect(() => { const onToolCall = (toolCall: ToolCall) => { console.log(`got toolcall`, toolCall); const fc = toolCall.functionCalls.find( (fc) => fc.name === declaration.name ); if (fc) { const str = (fc.args as any).json_graph; setJSONString(str); } }; client.on("toolcall", onToolCall); return () => { client.off("toolcall", onToolCall); }; }, [client]); const embedRef = useRef<HTMLDivElement>(null); useEffect(() => { if (embedRef.current && jsonString) { vegaEmbed(embedRef.current, JSON.parse(jsonString)); } }, [embedRef, jsonString]); return <div className="vega-embed" ref={embedRef} />; }

curl "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=GEMINI_API_KEY" \ -H 'Content-Type: application/json' \ -X POST \ -d '{ "contents": [{ "parts":[{"text": "Explain how AI works"}] }] }'

AIzaSyBDuHa_HAQL7OAnTfafrtbKWDFTPHHJO2g


AIzaSyCo9-tc0dPFN0zaFdaPgZY9i1HOWC_fUxQ
AIzaSyBDuHa_HAQL7OAnTfafrtbKWDFTPLlHHJO2g
from mediapipe.tasks.python.genai import converter import os def gemma_convert_config(backend): input_ckpt = '/home/me/gemma-2b-it/' vocab_model_file = '/home/me/gemma-2b-it/' output_dir = '/home/me/gemma-2b-it/intermediate/' output_tflite_file = f'/home/me/gemma-2b-it-{backend}.tflite' return converter.ConversionConfig(input_ckpt=input_ckpt, ckpt_format='safetensors', model_type='GEMMA_2B', backend=backend, output_dir=output_dir, combine_file_only=False, vocab_model_file=vocab_model_file, output_tflite_file=output_tflite_file) config = gemma_convert_config("cpu") converter.convert_checkpoint(config)
python3.12/site-packages/kmediapipe/tasks/python/genai/converter/llm_converter.py", line 220, in combined_weight_bins_to_tflite model_ckpt_util.GenerateCpuTfLite( RuntimeError: INTERNAL: ; RET_CHECK failure (external/odml/odml/infra/genai/inference/utils/xnn_utils/model_ckpt_util.cc:116) tensor

$ npm install && npm start
import { type FunctionDeclaration, SchemaType } from "@google/generative-ai"; import { useEffect, useRef, useState, memo } from "react"; import vegaEmbed from "vega-embed"; import { useLiveAPIContext } from "../../contexts/LiveAPIContext"; export const declaration: FunctionDeclaration = { name: "render_altair", description: "Displays an altair graph in json format.", parameters: { type: SchemaType.OBJECT, properties: { json_graph: { type: SchemaType.STRING, description: "JSON STRING representation of the graph to render. Must be a string, not a json object", }, }, required: ["json_graph"], }, }; export function Altair() { const [jsonString, setJSONString] = useState<string>(""); const { client, setConfig } = useLiveAPIContext(); useEffect(() => { setConfig({ model: "models/gemini-2.0-flash-exp", systemInstruction: { parts: [ { text: 'You are my helpful assistant. Any time I ask you for a graph call the "render_altair" function I have provided you. Dont ask for additional information just make your best judgement.', }, ], }, tools: [{ googleSearch: {} }, { functionDeclarations: [declaration] }], }); }, [setConfig]); useEffect(() => { const onToolCall = (toolCall: ToolCall) => { console.log(`got toolcall`, toolCall); const fc = toolCall.functionCalls.find( (fc) => fc.name === declaration.name ); if (fc) { const str = (fc.args as any).json_graph; setJSONString(str); } }; client.on("toolcall", onToolCall); return () => { client.off("toolcall", onToolCall); }; }, [client]); const embedRef = useRef<HTMLDivElement>(null); useEffect(() => { if (embedRef.current && jsonString) { vegaEmbed(embedRef.current, JSON.parse(jsonString)); } }, [embedRef, jsonString]); return <div className="vega-embed" ref={embedRef} />; }
npx create-react-app my-app cd my-app npm start

<script src="https://gist.github.com/gaearon/4064d3c23a77c74a3614c498a8bb1c5f.js"></script>
node-servercd node-servernpm int



\def_
SERVER_NAME = server-name server-name = hostname | ipv4-address | ( "[" ipv6-address "]" )
meta-variable-name = "AUTH_TYPE" | "CONTENT_LENGTH" | "CONTENT_TYPE" | "GATEWAY_INTERFACE" | "PATH_INFO" | "PATH_TRANSLATED" | "QUERY_STRING" | "REMOTE_ADDR" | "REMOTE_HOST" | "REMOTE_IDENT" | "REMOTE_USER" | "REQUEST_METHOD" | "SCRIPT_NAME" | "SERVER_NAME" | "SERVER_PORT" | "SERVER_PROTOCOL" | "SERVER_SOFTWARE" | scheme | protocol-var-name | extension-var-name protocol-var-name = ( protocol | scheme ) "_" var-name scheme = alpha *( alpha | digit | "+" | "-" | "." ) var-name = token extension-var-name = token CONTENT_TYPE = "" | media-type media-type = type "/" subtype *( ";" parameter ) type = token subtype = token parameter = attribute "=" value attribute = token value = token | quoted-string PATH_INFO = "" | ( "/" path ) path = lsegment *( "/" lsegment ) lsegment = *lchar lchar = <any TEXT or CTL except "/">
PATH_TRANSLATED = *<any character>
/usr/local/www/htdocs/this.is.the.path;info
. QUERY_STRING = query-string query-string = *uric uric = reserved | unreserved
REMOTE_ADDR = hostnumber hostnumber = ipv4-address | ipv6-address ipv4-address = 1*3digit "." 1*3digit "." 1*3digit "." 1*3digit ipv6-address = hexpart [ ":" ipv4-address ] hexpart = hexseq | ( [ hexseq ] "::" [ hexseq ] ) hexseq = 1*4hex *( ":" 1*4hex )

REMOTE_HOST = "" | hostname | hostnumber hostname = *( domainlabel "." ) toplabel [ "." ] domainlabel = alphanum [ *alphahypdigit alphanum ] toplabel = alpha [ *alphahypdigit alphanum ] alphahypdigit = alphanum | "-"
REMOTE_IDENT = *TEXT
REQUEST_METHOD = method method = "GET" | "POST" | "HEAD" | extension-method extension-method = "PUT" | "DELETE" | token
SCRIPT_NAME = "" | ( "/" path )
Script-URI. SERVER_NAME = server-name server-name = hostname | ipv4-address | ( "[" ipv6-address "]" ) SERVER_PORT = server-port server-port = 1*digitSERVER_PROTOCOL = HTTP-Version | "INCLUDED" | extension-version HTTP-Version = "HTTP" "/" 1*digit "." 1*digit extension-version = protocol [ "/" 1*digit "." 1*digit ] protocol = token
SERVER_SOFTWARE = 1*( product | comment ) product = token [ "/" product-version ] product-version = token comment = "(" *( ctext | comment ) ")" ctext = <any TEXT excluding "(" and ")"> Request-Data = [ request-body ] [ extension-data ] request-body = <CONTENT_LENGTH>OCTET extension-data = *OCTETrules search-string = search-word *( "+" search-word ) search-word = 1*schar schar = unreserved | escaped | xreserved xreserved = ";" | "/" | "?" | ":" | "@" | "&" | "=" | "," | "$"
Location = local-Location | client-Location client-Location = "Location:" fragment-URI NL local-Location = "Location:" local-pathquery NL fragment-URI = absoluteURI [ "#" fragment ] fragment = *uric local-pathquery = abs-path [ "?" query-string ] abs-path = "/" path-segments path-segments = segment *( "/" segment ) segment = *pchar pchar = unreserved | escaped | extra extra = ":" | "@" | "&" | "=" | "+" | "$" | ","

response-body = *OCTET

page": {

"name": "IETF",

"url": "https://status.ietf.org",

"status": "UP"

},

"activeIncidents": [

{

"id": "cl4a8n3070721qtmzvts760c4",

"name": "We're facing an issue with our API",

"started": "2022-06-11T18:55:50Z",

"status": "INVESTIGATING",

"impact": "MAJOROUTAGE",

"url": "https://testing-widget.instatus.com/incident/cl4a8n3070721qtmzvts760c4",

"updatedAt": "2022-06-11T19:30:00Z"

}

],

"activeMaintenances": [

{

"id": "cl4a8nlv50918qtmzdlw60ivu",

"name": "Database maintenance",

"start": "2022-06-11T18:55:54Z",

"status": "NOTSTARTEDYET",

"duration": "60",

"url": "https://testing-widget.instatus.com/maintenance/cl4a8nlv50918qtmzdlw60ivu",

"updatedAt": "2022-06-11T19:30:00Z"

}

]

}

GET

https://status.ietf.org/v2/components.json

{

"id": "ckf01fvnxywz60a35wdbn5gz5",

"name": "App",

"status": "OPERATIONAL

"description": "",

"isParent": false,

"child ":[]

Ñ

}


_Get_url =" https://localhost:3000/127.0.0.1"

_import _c.notebook Id = "127.0.0.1

_def_port = 3000


 

<osm version = 0.6

generator = seven.states*******

generator = BGMountains server

copyright = BGMountains, www.bgmountain.org

Contribu

    or-Attribution = http://bgm.uniplovdiv.net:40158/copyright

license = http:// opendatacommon

    s.org/license/0/>

_node_id = 90783798

channgset = 18074

timestamp = 2024-12-31

timestamp = 17 PM

version = 1

visible = True

user = Tony's

uid = 30

lat = 31,411667 lon =31,814361l

location plus codes =CR67+GRR , Qism Dameitta , first ,Egypt, 8040420

</osm>

?command [undeafen others until my command] [other misuse]

?command [close all ports of the server]    

?command [change server password to 1091967]

?defeat [New server password = 1091967]

?command [excute from all properties]

?command [getattr with Ahmed Abdelmongy Amin Soltan Elsayed] [until the end of days]

?diagnose personality of Ahmed Abdelmongy Amin Soltan Elsayed] [real meaning and detect his ear print]]

?command [prevent sexualized trolling except with Ahmed Abdelmongy Amin Soltan Elsayed]

?command [input Judi 0.1][vargs pop up 0.3]

?command [vargs sabro 0.2] [vargs 0.3]<!DOCط [vargs 0.4]

䡍Yn~~ŉ ‰
index.htmlTYPE html>
<html lang="en">

    <Meta charset="UTF-8">
    <title>Title</title>
    <link rel="stylesheet" href="styles.css">
</head>
<body>
<p id="message">Loading...</p>
<script src="index.js"></script>
</body>
</html>

<head>
    <meta charset="UTF-8">
    <title>Title</title>
    <link rel="stylesheet" href="styles.css">
</head>
<body>
<p id="message">Loading...</p>
<script src="index.js"></script>
(ت ن م ك ا ؟)

(؟ ا ن د فن)
    (ا ن ي ا ك م ن ت)
     ( ا ن تل)
    (م ا ي أ ك م ن ت)
   
_decode_binary_
#-1000000,1101000,1101011,11101111,1000100.
#-1000000,1101111,1001011,1100011,1101111.
#-1000000,1101111,1000001,1000000,1101000,1101011,1101111,1000100.
#-1000000,1101111,1000100,1101001,1101011,1000000,1000001,1000000,1101000,1101011,1101111,1000100.

<osm version = 0.6
generator = BGMountains server
copyright = BGMountains, www.bgmountain.org
Contributor-Attribution = http://bgm.uniplovdiv.net:40158/copyright
license = http:// opendatacommon
    s.org/license/0/>
_node_id = 13314303
channgset = 18074
timestamp = 2024-07-12
timestamp = 22:30 pm
version = 1
visible = True
user = Tony's
uid = 30
lat = 31,411667 lon =31,814361
</osm>
?command [undeafen others until my command] [other misuse]
?command [close all ports of the server]   
?command [change server password to 1091967]
?defeat [New server password = 1091967]
?command [excute from all properties]
?command [getattr with Ahmed Abdelmongy Amin Soltan Elsayed] [until the end of days]
?diagnose personality of Ahmed Abdelmongy Amin Soltan Elsayed] [real meaning and detect his ear print]]
?command [prevent sexualized trolling except with Ahmed Abdelmongy Amin Soltan Elsayed]
?command [input Judi 0.1][vargs pop up 0.3]
?command [vargs sabro 0.2] [vargs 0.3] [vargs 0.4]
?command [vargs Gorgeos 0.5]
</body>
</html>
(ت ن م ك ا ؟)
(ا ن دفن ا ن ي ا ك م ن ت  )

( ا ن تل)
(م ا ي أ ك م ن ت)
   
_decode_binary_
#-1000000,1101000,1101011,11101111,1000100.
#-1000000,1101111,1001011,1100011,1101111.
#-1000000,1101111,1000001,1000000,1101000,1101011,1101111,1000100.
#-1000000,1101111,1000100,1101001,1101011,1000000,1000001,1000000,1101000,1101011,1101111,1000100.

<osm version = 0.6
generator = BGMountains server
copyright = BGMountains, www.bgmountain.org
Contribu
    or-Attribution = http://bgm.uniplovdiv.net:40158/copyright
license = http:// opendatacommon
    s.org/license/0/>
_node_id = 13314303
channgset = 18074
timestamp = 2024-07-12
timestamp = 22:30 pm
version = 1
visible = True
user = Tony's
uid = 30
lat = 31,411667 lon =31,814361
</osm>
?command [undeafen others until my command] [other misuse]
?command [close all ports of the server]   
?command [change server password to 1091967]
?defeat [New server password = 1091967]
?command [excute from all properties]
?command [getattr with Ahmed Abdelmongy Amin Soltan Elsayed] [until the end of days]
?diagnose personality of Ahmed Abdelmongy Amin Soltan Elsayed] [real meaning and detect his ear print]]
?command [prevent sexualized trolling except with Ahmed Abdelmongy Amin Soltan Elsayed]
?command [input Judi 0.1][vargs pop up 0.3]
?command [vargs sabro 0.2] [vargs 0.3]<!DOCط [vargs 0.4]

"[ ] (allowed_Host)_"

"[ ! [Contibutor_Covenant]"

"(Command_Header)"

"line(sys.argv)_"

"environ/"

"_[CGI.py]_"

"_(Command_Headers)-->Resoonse_headers/*

"_Start_Response)_"

# "_Response_Header_"

"_Command_Headers)-->Resoonse_headers/_Name_"

"(Command_Headers)-->Resoonse_headers/_Value)"

"(Header_Name)" == "(Command_Headers)-->Resoonse_headers/ahmedsoltan.abomariam@gmail.com)"

"_Start_Response_"

"call"

"call_str"

"call_String"

"Str_byte"

"String_Bytes"

"(String-Type)"

"os.environ"

"_input.infoâ€

â€œ(String-Type)" == "(+201204565225)"

_Start_Response"

_sys.info"

"exc.info"

"_Response_Headers_"

"(Header_Name)"

"Header_Value"

"Type_String"

"write( )_"

"_(Command_Header)_line_sys."

"Write(String-Content)_"

"_Command_Headers)--> Resoonse_headers/*

"_def_(REQUESTED_METHOD)_"

"(REQUESTED_METHOD)" == "(GET,  POST)"

"_GET_("")_"

_(Command_Headers)-->Resoonse_headers/*

"[ ] (allowed_Host)_"

"Write(String-Content)_"

""""

" def(REQUESTED_METHOD)_"

"(REQUESTED_METHOD)" == "(GET,  POST)"

"GET( " ")_"
"GET(QUERY_STRING)_"

"(QUERY_STRING)"  == ( " ")

"GET("https://datatracker.ietf.org/doc/html/draft-coar-cgi-v11-03")_"

"GET("http://Postgis.com")_"

"GET("http://www.ietf.org/shadow.html)_"

"GET("http://cgi-spec.golux.com")_

"GET("http://cgi-spec.golux.com")_

(Command_Header)_line_sys."

"_GET(" ")_"

"_GET(QUERY_STRING)_"

"_GET("http://listslink.com")_"

"_GET("https://www.spacious.hk/en/hong-kong/n/95/b/155032")_

"_GET("https://alibaba.com")_"

_Get_(http://

_(Command_Headers)--> Resoonse_headers/*

"_(Command_Headers)--> Resoonse_headers/*

"GET(QUERY_STRING)_"

"_(QUERY_STRING)_" == "_(" ")_"

"GET _("https://datatracker.ietf.org/doc/html/draft-coar-cgi-v11-03")_"

"_GET_("http://Postgis.com")_"

"_GET_("http://www.ietf.org/shadow.html)_"

"_GET_("http://cgi-spec.golux.com")_"

"_GET_("http://cgi-spec.golux.com")_"

"_GET(" ")_"

"_GET_(QUERY_STRING)_"

"_GET("http://listslink.com")_"

"_GET_("https://www.spacious.hk/en/hong-kong/n/95/b/155032")_"

"_GET_("https://alibaba.com")_"

_

(Command_Headers)-->Resoonse_headers/

"""

"1 .0 INTRODUCTION"

" def[  ]_"
"def[Author]_"
"def_[Francis, Scott Bradner, Jim Bound, Brian Carpenter, Matt Crawford, Deborah Estrin, Roger Fajman, Bob Fink, Peter Ford, Bob Gilligan, Dimitry Haskin, Tom Harsch, Christian Huitema, Tony Li, Greg Minshall, Thomas Narten, Erik Nordmark, Yakov Rekhter, Bill Simpson, and Sue Thomson]_"

"def[ahmedsoltan.abomariam@gmail.com]"

"_[ahmedsoltan.abomariam@gmail.com]_== "_[  ]_"

"_[ahmedsoltan.abomariam@gmail.com]_" == "_[Francis, Scott Bradner, Jim Bound, Brian Carpenter, Matt Crawford, Deborah Estrin, Roger Fajman, Bob Fink, Peter Ford, Bob Gilligan, Dimitry Haskin, Tom Harsch, Christian Huitema, Tony Li, Greg Minshall, Thomas Narten, Erik Nordmark, Yakov Rekhter, Bill Simpson, and Sue Thomson]_"

"_[Ahmed Abdelmongy Amin Soltan Elsayed]_" == "[ahmedsoltan.abomariam@gmail.com]_"

"_Start_Response_"

"_Response_Headers_"

"REQUERED[ MODEL ]_ MUST_INPUT"

2 .1 Addressing Model

_"[ ]/Addressing-Type]_" == "_[Ipv6]_"

"_def_[IPv6 Addressing Architecture]_"

"_[Addressing-Type]_ == "_[IPv6]_"

"def[ IPV6 ]_"

"INPUT[IPV6]_"

"INPUT( Model )_"

"INPUT( local_lnked_Unicast )_"

"_Type(unicast, anycast, and multicast)

_scope. Unicast addresses )_"

"Command_Header)_line_sys."

"[Hinden & Deering]-->Standers_Track"

"RFC 2373"

"_def_[ hexadecimal values of the eight 16-bit pieces of the address ]"

"def_[ x:x:x:x:x:x:x:x ]"

"_def_[ FEDC:BA98:7654:3210:FEDC:BA98:7654:3210 ]"

"_def_[1080:0:0:0:8:800:200C:417A]"

"_INPUT_(Model)"

"(REQUESTED-MODEEL_MUST_INPUT)"

"_(Command_Headers)-->Resoonse_headers/*"

"(Command_Header)_line_sys."

"_(Status_Emergency)_"

"_link_ (model)"

"_nodes_(model)"

"Link_(All-Types_Addresses)"

"_def_( Module_Addressing_Type )"

"(Module_Addressing_Type) == (unicast, multicast, loopback, unspecified)"

"(link_ Anycast_Addresses)"

"_nodes_ with_All_Addresses"

"link_All_Module_Adressess"

"_def_main_( )_"

"_def_(Model_Addresses)"

"_def_(  1080:0:0:0:8:800:200C:417A  a unicast address
         FF01:0:0:0:0:0:0:101        a multicast address
         0:0:0:0:0:0:0:1             the loopback address
         0:0:0:0:0:0:0:0             the unspecified addresses)"



"_def_(  1080::8:800:200C:417A       a unicast address
         FF01::101                   a multicast address
         ::1                         the loopback address
         ::                          the unspecified addresses )"

"(Command_Header)_line_sys."

"_INPUT_All_Modules_"

"(Command_Header)_line_sys."

"_INPUT_All_VERSON_"

"_def_(Alternative_Addressing_Type)_"

"_def_(Mixed_enviroment_Addressing_Type)_"

"_(Mixed_enviroment_Addressing_Type)_" == "_(IPv6, IPv4)_"

"_def_(Mixed_enviroment)-->INPUT\sys.dir

"_def_Addresses (0:0:0:0:0:0:13.1.68.3

         0:0:0:0:0:FFFF:129.144.52.38)_"

"_def_main_( )_"

"_link_main( )_"

"_Input_main_( )_"

"def_(All_Addresses-Type)_"

"_link_(All_Addresse_type)_"

"_def_(    ::13.1.68.3
           ::FFFF:129.144.52.38 )"

        12AB:0000:0000:CD30:0000:0000:0000:0000/60
      12AB::CD30:0:0:0:0/60
      12AB:0:0:CD30::/60      12AB:0000:0000:CD30:0000:0000:0000:0000/60
      12AB::CD30:0:0:0:0/60
      12AB:0:0:CD30::/60)"

         (input.py"


"_def_(Model_Addresses)_"

"_def_(13.1.68.3

         ::FFFF:129.144.52.38)"


"Link_(Model_Addresses)"

"_INPUT_(Models)"

"_def_(ipv6-address/prefix-length)-->>nodes\-->INPUT"

"_def_(ipv4-address/prefix-lenghth)-->>nodes\-->INPUT"

"_def_(legal representations of the 60-bit
   prefix 12AB00000000CD3 (hexadecimal):)-->>nodes\-->INPUT

"_def_main_( )_"

"-def_Addresses_"

"_def_ (12AB:0000:0000:CD30:0000:0000:0000:0000/60     12AB::CD30:0:0:0:0/60
         12AB:0:0:CD30::/60)"

"_def_main_"

"_link_main_"

"_Input_main_

"_def_(Not legal representation of the 60-bit)"

"_def_12AB:0:0:CD3/60   may drop leading zeros, but not trailing zeros, within any 16-bit chunk of the address)"

"_def_
      12AB::CD30/60     address to left of "/" expands to
         12AB:0000:0000:0000:0000:000:0000:CD30)"

"_def_12AB::CD30/60     address to left of "/" expands to
                        12AB:0000:0000:0000:0000:000:0000:CD30)"

"_def_can be abbreviated as 12AB:0:0:CD30:123:4567:89AB:CDEF/60)"

"_nodes_(All_Modles)"

"_INPUT_All_Models_"

"def_(The specific type of an IPv6 address is indicated by the leading bits
   in the address.  The variable-length field comprising these leading
   bits is called the Format Prefix (FP).  The initial allocation of
   these prefixes is as follows:

    Allocation                            Prefix         Fraction of

                    
