function is_array(obj){
  //проверка - является ли obj - массивом
    return obj !== null && obj instanceof Array;
}
function is_object(obj){
    return obj !== null && typeof obj == 'object';
}
function var_dump(obj,max_level,level,separator){
  if(max_level == null) max_level = 99999;
  if(level == null) level = 0;
  if(separator==null) separator = '\n';
  if(is_object(obj)){
    if(level>=max_level) return "{/*dump_max_level:"+max_level+"*/}";
    var tab = '  ';
    var tab_lv = tab;
    var i = 0;
    while(i++<level) tab_lv += tab;
    var ret = '{\n';
    //var end = obj.length;
    for (var key in obj){ // обращение к свойствам объекта по индексу
      ret += tab_lv + key + " : " + var_dump(obj[key],max_level,level+1) ;
      ret += ',';
      ret += separator;
    }
    ret += tab_lv + '}';
    return ret;
  }
  return obj;
}
//----------------------------------------------------

$( document ).ready(function() {
	loaddatahostfileslist();
	loaddataukfileslist();
	loadversionlist();
	
	$("#template_ver").on("change",loadfileslist);
	$("#template_file").on("change",loadfiledata);
	$("#render_query").on("click",renderquery);
	$("#send_query").on("click",sendquery);
})

//загрузка списка файлов данных хостов
function loaddatahostfileslist(){
	$.ajax({
		type: "GET",
		url: "/loaddatahostfileslist",
		data: {},
		dataType: "json",
		beforeSend: function(){
			$("#result_xml").html("<option selected>загрузка...</option>")
		}
	}).done(function(json){
		if(json.err) return show_error_ajax_fail(json)
		var s = "";
		n = 0;
		for(var i in json){
			if(n++==json.length-1) s += "<option selected>"+json[i]+"</option>"
			else s += "<option>"+json[i]+"</option>"
		}
		$("#data_host").html(s)
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail(xhr, status, errorThrown)
  	});
}


//загрузка списка версий
function loadversionlist(){
	$.ajax({
		type: "GET",
		url: "/loadversionlist",
		dataType: "json",
		beforeSend: function(){
			$("#template_ver").html("<option selected>загрузка...</option>")
		}
	}).done(function(json){
		if(json.err) return show_error_ajax_fail(json)
		var s = "";
		n = 0;
		for(var i in json){
			if(n++==json.length-1) s += "<option selected>"+json[i]+"</option>"
			else s += "<option>"+json[i]+"</option>"
		}
		$("#template_ver").html(s)
		loadfileslist()
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail(xhr, status, errorThrown)
  	});
}

//загрузка списка xml файлов в версии
function loadfileslist(){
	$.ajax({
		type: "GET",
		url: "/loadfileslist",
		data: {ver: $("#template_ver").val()},
		dataType: "json",
		beforeSend: function(){
			$("#template_file").html("<option selected>загрузка...</option>")
		}
	}).done(function(json){
		if(json.err) return show_error_ajax_fail(json)
		var s = "";
		for(var i in json){
			s += "<option>"+json[i]+"</option>"
		}
		$("#template_file").html(s)
		loadfiledata()
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail(xhr, status, errorThrown)
  	});
}

//загрузка списка файлов данных по ук
function loaddataukfileslist(){
	$.ajax({
		type: "GET",
		url: "/loaddataukfileslist",
		data: {},
		dataType: "json",
		beforeSend: function(){
			$("#result_xml").html("<option selected>загрузка...</option>")
		}
	}).done(function(json){
		if(json.err) return show_error_ajax_fail(json)
		var s = "";
		n = 0;
		for(var i in json){
			if(n++==json.length-1) s += "<option selected>"+json[i]+"</option>"
			else s += "<option>"+json[i]+"</option>"
		}
		$("#data_uk").html(s)
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail(xhr, status, errorThrown)
  	});
}

//загрузка данных по выбранному файлу
function loadfiledata(){
	$.ajax({
		type: "GET",
		url: "/loadfiledata",
		data: {ver: $("#template_ver").val(), filename: $("#template_file").val()},
		dataType: "json",
		beforeSend: function(){
			$("#template_data").val("загрузка...")
			$("#template_xml").val("загрузка...")
		}
	}).done(function(json){
			if(json.err) return show_error_ajax_fail(json)
			$("#template_data").val(json[1])
			$("#template_xml").val(json[0])
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail(xhr, status, errorThrown)
  	});
}

//генерация шаблона по занным параметрам
function renderquery(){
	var data = {}
	data.ver = $("#template_ver").val()
	data.datafilename_host = $("#data_host").val()
	data.datafilename_uk = $("#data_uk").val()
	data.data = $("#template_data").val()
	data.xml = $("#template_xml").val()
	
	$.ajax({
		type: "POST",
		url: "/renderquery",
		data: data,
		dataType: "json",
		beforeSend: function(){
			$("#render_data").val("загрузка...")
			$("#render_xml").val("загрузка...")
		}
	}).done(function(json){
			if(json.err) return show_error_ajax_fail(json)
			$("#render_data").val(json[1])
			$("#render_xml").val(json[0])
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail(xhr, status, errorThrown)
  	});
}

//отправка запроса
function sendquery(){
	var data = {}
	data.data = $("#render_data").val()
	data.xml = $("#render_xml").val()
	
	$.ajax({
		type: "POST",
		url: "/sendquery",
		data: data,
		dataType: "json",
		beforeSend: function(){
			$("#result_data").val("загрузка...")
			$("#result_xml").val("загрузка...")
		}
	}).done(function(json){
			if(json.err) return show_error_ajax_fail(json)
			$("#result_data").val(json[1])
			$("#result_xml").val(json[0])
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail(xhr, status, errorThrown)
  	});
}

function show_error_ajax_fail(xhr, status, errorThrown){
	if(arguments.length==1){
		json = xhr;
		$("#errinfo").html(
		  '<div class="modal" id="errModal" role="dialog"><div class="modal-dialog">                '+
		  '    <div class="modal-content"><div class="modal-header">                                '+
		  '        <button type="button" class="close" data-dismiss="modal">&times;</button>        '+
		  '        <h4 class="modal-title">req error</h4>                                           '+
		  '      </div>                                                                             '+
		  '      <div class="modal-body">                                                           '+
		  '        <pre>'+var_dump(json)+'</pre>                                                    '+
		  '      </div>                                                                             '+
		  '      <div class="modal-footer">                                                         '+
		  '        <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>'+
		  '      </div>                                                                             '+
		  '    </div>                                                                               '+
		  '</div></div>'
		);
		$("#errinfo").find("#errModal").modal()
		return
	}
	$("#errinfo").html(
	  '<div class="modal" id="errModal" role="dialog"><div class="modal-dialog">                '+
	  '    <div class="modal-content"><div class="modal-header">                                '+
	  '        <button type="button" class="close" data-dismiss="modal">&times;</button>        '+
	  '        <h4 class="modal-title">'+var_dump(status)+'</h4>                                '+
	  '      </div>                                                                             '+
	  '      <div class="modal-body">                                                           '+
	  '        <pre>'+var_dump(xhr.responseText)+'</pre>                                        '+
	  '      </div>                                                                             '+
	  '      <div class="modal-footer">                                                         '+
	  '        <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>'+
	  '      </div>                                                                             '+
	  '    </div>                                                                               '+
	  '</div></div>'
	);
	$("#errinfo").find("#errModal").modal()
}

