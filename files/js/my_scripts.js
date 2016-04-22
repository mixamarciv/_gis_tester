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
	if($("#main_form").length>0) {
		
		$("#savexmlfile").on("click",savexmlfile);
		$("#delxmlfile").on("click",delxmlfile);

		loaddatahostfileslist();
		loaddataukfileslist();
		loadversionlist();
		
		loadresultlist();
		
		$("#template_ver").on("change",loadfileslist);
		$("#template_file").on("change",loadfiledata);
		$("#render_query").on("click",renderquery);
		$("#resutsfilter").on("change",loadresultlist);
		$("#resutsfilter").on("input",loadresultlist);
		$("#load_result").on("click",loadresult);
		
		
		$("#send_query").on("click",sendquery);
		$("#save_result").on("click",saveresult);
		
		
		$("#data_uk_edit").on("click",function(){ editfile("data_uk") });
		$("#data_host_edit").on("click",function(){ editfile("data_host") });
		//type == data_uk или data_host
	}
	
	if($("#edit_form").length>0) {
		$("#save").on("click",edit_form_save);
		$("#delete").on("click",edit_form_delete);
	}
})

function edit_form_save(){
	var data = {}
	data.data = $("#data").val()
	//data.json = ''
	data.file = $("#file").val()
	data.rtype = $("#rtype").val()
	data.oper = "save"
	
	$.ajax({
		type: "POST",
		url: "/edit_save",
		data: data,
		dataType: "json",
		beforeSend: function(){
			$("#info").html("загрузка...")
		}
	}).done(function(json){
			if(json.err) return show_error_ajax_fail("renderquery",json)
			$("#info").html("файл успешно сохранен "+json.file)
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail("renderquery",xhr, status, errorThrown)
  	});
}

function edit_form_delete(){
	var data = {}
	data.file = $("#file").val()
	data.rtype = $("#rtype").val()
	data.oper = "delete"
	
	$.ajax({
		type: "POST",
		url: "/edit_save",
		data: data,
		dataType: "json",
		beforeSend: function(){
			$("#info").html("загрузка...")
		}
	}).done(function(json){
			if(json.err) return show_error_ajax_fail("renderquery",json)
			$("#info").html("файл успешно удален "+json.file)
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail("renderquery",xhr, status, errorThrown)
  	});
}


function savexmlfile(){
	var data = {}
	data.json = $("#template_data").val()
	data.xml = $("#template_xml").val()
	data.file = $("#savexmlfilename").val()
	data.rtype = "xml/"+$("#template_ver").val()
	data.oper = "savexml"
	
	
	$.ajax({
		type: "POST",
		url: "/edit_save",
		data: data,
		dataType: "json",
		beforeSend: function(){
			$("#savexmlinfo").html("загрузка...")
		}
	}).done(function(json){
			if(json.err) return show_error_ajax_fail("renderquery",json)
			$("#savexmlinfo").html("сохранено")
			loadfileslist(data.file);
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail("renderquery",xhr, status, errorThrown)
  	});
}

function delxmlfile(){
	var data = {}
	data.file = $("#savexmlfilename").val()
	data.rtype = "xml/"+$("#template_ver").val()
	data.oper = "delxml"
	
	if (confirm("удалить файл "+data.file+"?")){
		$.ajax({
			type: "POST",
			url: "/edit_save",
			data: data,
			dataType: "json",
			beforeSend: function(){
				$("#savexmlinfo").html("загрузка...")
			}
		}).done(function(json){
				if(json.err) return show_error_ajax_fail("renderquery",json)
				$("#savexmlinfo").html("удалено")
				loadfileslist(data.file);
		}).fail(function( xhr, status, errorThrown ) {
			show_error_ajax_fail("renderquery",xhr, status, errorThrown)
	  	});
	}
}

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
		if(json.err) return show_error_ajax_fail("loaddatahostfileslist",json)
		var s = "";
		n = 0;
		for(var i in json){
			if(n++==json.length-1) s += "<option selected>"+json[i]+"</option>"
			else s += "<option>"+json[i]+"</option>"
		}
		$("#data_host").html(s)
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail("loaddatahostfileslist",xhr, status, errorThrown)
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
		if(json.err) return show_error_ajax_fail("loadversionlist",json)
		var s = "";
		n = 0;
		for(var i in json){
			if(n++==json.length-1) s += "<option selected>"+json[i]+"</option>"
			else s += "<option>"+json[i]+"</option>"
		}
		$("#template_ver").html(s)
		loadfileslist()
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail("loadversionlist",xhr, status, errorThrown)
  	});
}

//загрузка списка xml файлов в версии
function loadfileslist(file){
	$.ajax({
		type: "GET",
		url: "/loadfileslist",
		data: {ver: $("#template_ver").val()},
		dataType: "json",
		beforeSend: function(){
			$("#template_file").html("<option selected>загрузка...</option>")
		}
	}).done(function(json){
		if(json.err) return show_error_ajax_fail("loadfileslist",json)
		var s = "";
		for(var i in json){
			sel = ""
			if(file){ 
				if(file==json[i]) sel = "selected"
			}else if(n++==json.length-1) sel = "selected"
			s += "<option "+sel+">"+json[i]+"</option>"
		}
		$("#template_file").html(s)
		loadfiledata()
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail("loadfileslist",xhr, status, errorThrown)
  	});
}

//загрузка списка файлов данных по ук
function loaddataukfileslist(file){
	$.ajax({
		type: "GET",
		url: "/loaddataukfileslist",
		data: {},
		dataType: "json",
		beforeSend: function(){
			$("#result_xml").html("<option selected>загрузка...</option>")
		}
	}).done(function(json){
		if(json.err) return show_error_ajax_fail("loaddataukfileslist",json)
		var s = "";
		n = 0;
		for(var i in json){
			sel = ""
			if(file){ 
				if(file==json[i]) sel = "selected"
			}else if(n++==json.length-1) sel = "selected"
			s += "<option "+sel+">"+json[i]+"</option>"
		}
		$("#data_uk").html(s)
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail("loaddataukfileslist",xhr, status, errorThrown)
  	});
}

//загрузка данных по выбранному файлу
function loadfiledata(){
	$("#savexmlfilename").val($("#template_file").val())
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
			if(json.err) return show_error_ajax_fail("loadfiledata",json)
			$("#template_data").val(json[1])
			$("#template_xml").val(json[0])
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail("loadfiledata",xhr, status, errorThrown)
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
			if(json.err) return show_error_ajax_fail("renderquery",json)
			$("#render_data").val(json[1])
			$("#render_xml").val(json[0])
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail("renderquery",xhr, status, errorThrown)
  	});
}

//отправка запроса (сам запрос и результат сохраняем в window.data1 что бы потом можно было им воспользоваться)
function sendquery(){
	var data = {}
	data.data = $("#render_data").val()
	data.xml = $("#render_xml").val()
	window.data1 = data;
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
			if(json.err) return show_error_ajax_fail("sendquery",json)
			$("#result_data").val(json[1])
			$("#result_xml").val(json[0])
			window.data1.res_data = json[1]
			window.data1.res_xml = json[0]
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail("sendquery",xhr, status, errorThrown)
  	});
}

function saveresult(){
	if(!window.data1){
		window.data1 = {}
		window.data1.data = $("#render_data").val()
		window.data1.xml  = $("#render_xml").val()
		window.data1.res_data = $("#result_data").val()
		window.data1.res_xml  = $("#result_xml").val()
	}
	window.data1.name_prefix = $("#save_result_fileprefix").val()
	//alert(window.data1.name_prefix)
	$.ajax({
		type: "POST",
		url: "/saveresult",
		data: window.data1,
		dataType: "json",
		beforeSend: function(){
		}
	}).done(function(json){
			if(json.err) return show_error_ajax_fail("saveresult",json)
			if(!json.ok) return show_error_ajax_fail("saveresult",json)
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail("saveresult",xhr, status, errorThrown)
  	});
}

function loadresultlist(){
	$.ajax({
		type: "POST",
		url: "/loadresultlist",
		data: {filter:$("#resutsfilter").val()},
		dataType: "json",
		beforeSend: function(){
			$("#resultslist").html("<option selected>load..</option>")
		}
	}).done(function(json){
		var s = "";
		if(!json){ 
			s="<option selected>not found</option>"
			$("#resultslist").html(s)
			$("#result_label").html("results(0)")
			return
		}
		if(json.err) return show_error_ajax_fail("loadresultlist",json)
		n = 0;
		for(var i in json){
			if(n++==json.length-1) s += "<option selected>"+json[i]+"</option>"
			else s += "<option>"+json[i]+"</option>"
			if(n>100) break;
		}
		if(n==0) s="<option selected>not found</option>"
		$("#resultslist").html(s)
		$("#result_label").html("results("+n+")")
	}).fail(function( xhr, status, errorThrown ) {
		//alert("ERR:"+var_dump(xhr))
		show_error_ajax_fail("loadresultlist",xhr, status, errorThrown)
  	});
}

//загрузка ранее сохраненных данных
function loadresult(){
	var data = {}
	$.ajax({
		type: "POST",
		url: "/loadresult",
		data: {file:$("#resultslist").val()},
		dataType: "json",
		beforeSend: function(){
			$("#render_data").val("загрузка...")
			$("#render_xml").val("загрузка...")
			$("#result_data").val("загрузка...")
			$("#result_xml").val("загрузка...")
		}
	}).done(function(json){
			//alert(var_dump(json))
			if(json.err) return show_error_ajax_fail("loadresult",json)
			$("#render_data").val(json.data)
			$("#render_xml").val(json.xml)
			$("#result_data").val(json.res_data)
			$("#result_xml").val(json.res_xml)
	}).fail(function( xhr, status, errorThrown ) {
		show_error_ajax_fail("loadresult",xhr, status, errorThrown)
  	});
}

//редактирование файлов для data_uk или data_host
function editfile(type){
	//type == data_uk или data_host
	var file = $("#"+type).val()
	var url = "/edit?type="+type+"&file="+encodeURIComponent(file)
	window.open(url,'_blank');
	/********
	window.open(url,'_blank');
	var newWin = window.open("/edit?type="+type+"&file="+encodeURIComponent(file),
	   "edit "+type,
	   "width=420,height=230,resizable=yes,scrollbars=yes,status=yes"
	)
	
	
	newWin.focus()
	*********/
}

function show_error_ajax_fail(info,xhr, status, errorThrown){
	if(arguments.length==2){
		json = xhr;
		$("#errinfo").html(
		  '<div class="modal" id="errModal" role="dialog"><div class="modal-dialog">                '+
		  '    <div class="modal-content"><div class="modal-header">                                '+
		  '        <button type="button" class="close" data-dismiss="modal">&times;</button>        '+
		  '        <h4 class="modal-title">'+info+'</h4>                                           '+
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
	  '        <h4 class="modal-title">'+info+'  '+var_dump(status)+'</h4>                      '+
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

