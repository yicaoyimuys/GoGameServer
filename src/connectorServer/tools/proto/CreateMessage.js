/**
 * Created by yangsong on 16/1/24.
 */
var fs = require("fs");
var program = require('./program/program.js');

var protoName = program.protoFile;
var projectName = program.projectName;

var protoFile = "./src/"+projectName+"/proto/"+protoName+".json";
var msgFilePath = "./src/"+projectName+"/proto/msg/";

var msgTemplate = fs.readFileSync("./src/"+projectName+"/tools/proto/template/msgTemplate.txt","utf-8");
var msgIdTemplate = fs.readFileSync("./src/"+projectName+"/tools/proto/template/msgIdTemplate.txt","utf-8");



var EncodeObj = {
    "uint8": "SetUint8",
    "int8": "SetInt8",
    "byte": "SetUint8",
    "uint16": "SetUint16",
    "int16": "SetInt16",
    "ushort": "SetUint16",
    "short": "SetInt16",
    "uint32": "SetUint32",
    "int32": "SetInt32",
    "uint64": "SetUint64",
    "int64": "SetInt64",
    "float": "SetFloat",
    "string": "SetString",
    "buffer": "SetBuffer"
};

var DecodeObj = {
    "uint8": "GetUint8",
    "int8": "GetInt8",
    "byte": "GetUint8",
    "uint16": "GetUint16",
    "int16": "GetInt16",
    "ushort": "GetUint16",
    "short": "GetInt16",
    "uint32": "GetUint32",
    "int32": "GetInt32",
    "uint64": "GetUint64",
    "int64": "GetInt64",
    "float": "GetFloat",
    "string": "GetString",
    "buffer": "GetBuffer"
};

var PropertyObj = {
    "uint8": "uint8",
    "int8": "int8",
    "byte": "uint8",
    "uint16": "uint16",
    "int16": "int16",
    "ushort": "uint16",
    "float": "float32",
    "short": "int16",
    "uint32": "uint32",
    "int32": "int32",
    "uint64": "uint64",
    "int64": "int64",
    "string": "string",
    "buffer": "[]byte"
};

var imports = {};
var msgIdStr = {};
var idsStr = '';
var classNameStr = '';

buildFile();
generateMsgIdFile();

function replaceAll(str, s1, s2) {
    var demo = str.replace(s1, s2);
    while (demo.indexOf(s1) != - 1)
        demo = demo.replace(s1, s2);
    return demo;
}

function buildFile(){
    var msgObj = JSON.parse(readProtoFile());
    for (var key in msgObj){
        generateMsgFile(key, msgObj[key]);
    }
}

function readProtoFile(){
    var str = fs.readFileSync(protoFile).toString();
    str = str.replace(/\/\/.*[\n\r]/g, "");
    return str;
}

function generateMsgIdFile(){
    var fileContent = replaceAll(msgIdTemplate, "$0", idsStr);
    fileContent = replaceAll(fileContent, "$1", classNameStr);
    fileContent = replaceAll(fileContent, "$2", '"'+projectName+'/proto"');
    saveFile(program.protoFile+'MsgId', fileContent);
}

function generateMsgFile(fileName, msgObj){
    fileName = ucfirst(fileName);
    classNameStr += '\t' + 'proto.SetMsgByName("'+fileName+'", '+fileName+'{})' + '\n';
    if (msgObj.msgId) {
        classNameStr += '\t' + 'proto.SetMsgById('+msgObj.msgId+', '+fileName+'{})' + '\n';
    }

    imports[fileName] = [projectName+'/proto', 'bytes'];
    var fileContent = replaceAll(msgTemplate, "$1", fileName);
    var propertyStr = getPropertyStr(fileName, msgObj);
    var encodeStr = getEncodeStr(msgObj);
    var decodeStr = getDecodeStr(msgObj);
    var importStr = getImportStr(fileName);

    fileContent = fileContent.replace("$0", importStr);
    fileContent = fileContent.replace("$2", propertyStr);
    fileContent = fileContent.replace("$3", encodeStr);
    fileContent = fileContent.replace("$4", decodeStr);
    fileContent = fileContent.replace("$5", msgIdStr[fileName] || '')

    saveFile(fileName, fileContent);
}

function getImportStr(fileName){
    var str = '';
    imports[fileName].forEach(function(tmp){
        str += 'import "'+tmp+'"'
        str += '\n'
    })
    str += '\n'
    return str
}

function getColumnProperty(key, value){
    var str = ucfirst(key) + '\t';
    if(key == 'msgId'){
        str += 'uint16';
    }
    else if(value.indexOf('array') != -1){
        str += '*list.List'
    }
    else if(PropertyObj[value]){
        str += PropertyObj[value];
    }
    else {
        str += '*'+ucfirst(value);
    }
    return str;
}


function getPropertyStr(fileName, msgObj){
    var str = '';
    for(var key in msgObj){
        var value = msgObj[key];
        str += '\t' + getColumnProperty(key, value) + '\n'

        if(key == 'msgId'){
            idsStr += '\t' + 'ID_'+fileName+' uint16 = '+value+';' + '\n'
            msgIdStr[fileName] = 'MsgId: 	ID_'+fileName+','
        }

        if(value.indexOf('array') != -1){
            if(imports[fileName].indexOf('container/list') == -1){
                imports[fileName].push('container/list');
            }
        }
    }
    return str;
}

function ucfirst(word) {
    return word.substring(0, 1).toUpperCase() + word.substring(1);
}

function getColumnEncode(key, value){
    var str = '\t';
    if(key == 'msgId'){
        str += 'proto.SetUint16(buf, this.'+ucfirst(key)+')';
    }
    else if(value.indexOf('array') != -1){
        var arr = value.split('/');
        var arrType = arr[1];
        if(!EncodeObj[arrType]){
            arrType = ucfirst(arrType);
        }
        str += 'proto.SetArray(buf, this.'+ucfirst(key)+', "'+arrType+'")';
    }
    else if(EncodeObj[value]){
        str += 'proto.'+EncodeObj[value]+'(buf, this.'+ucfirst(key)+')';
    }
    else{
        str += 'proto.SetEntity(buf, this.'+ucfirst(key)+')';
    }
    str += '\n';
    return str;
}

function getEncodeStr(msgObj){
    var str = '';
    for(var key in msgObj){
        var value = msgObj[key];
        str += getColumnEncode(key, value);
    }
    return str;
}

function getColumnDecode(key, value){
    var str = '\t'
    if(key == 'msgId'){
        str += 'this.'+ucfirst(key)+' = proto.GetUint16(buf)';
    }
    else if(value.indexOf('array') != -1){
        var arr = value.split('/');
        var arrType = arr[1];
        if(!DecodeObj[arrType]){
            arrType = ucfirst(arrType);
        }
        str += 'this.'+ucfirst(key)+' = proto.GetArray(buf, "'+arrType+'")';
    }
    else if(DecodeObj[value]){
        str += 'this.'+ucfirst(key)+' = proto.'+DecodeObj[value]+'(buf)';
    }
    else{
        str += 'this.'+ucfirst(key)+' = proto.GetEntity(buf, "'+ucfirst(value)+'").(*'+ucfirst(value)+')';
    }
    str += '\n';
    return str;
}

function getDecodeStr(msgObj){
    var str = '';
    for(var key in msgObj){
        var value = msgObj[key];
        str += getColumnDecode(key, value);
    }
    return str;
}

function saveFile(fileName, content){
    fs.writeFileSync(msgFilePath + fileName + ".go", content);
}