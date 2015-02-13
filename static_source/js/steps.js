$(document).ready(function() {
	var i = 0;
	//	installProcess(i);

	var servers = ["MEGAM", "COBBLER"];
	install_check(i, servers);

	$('#MEGAM_install_button').click(function(event) {
		event.preventDefault();
		installProcess("MEGAM");
		return false;
		// for good measure
	});

	$('#COBBLER_install_button').click(function(event) {
		event.preventDefault();
		installProcess("COBBLER");
		return false;
		// for good measure
	});

	$('#OPENNEBULA_install_button').click(function(event) {
		event.preventDefault();
		installProcess("OPENNEBULA");
		return false;
		// for good measure
	});

	$('#OPENNEBULAHOST_install_button').click(function(event) {
		event.preventDefault();
		installProcess("OPENNEBULAHOST");
		return false;
		// for good measure
	});
	
	$('#COMPUTE_install_button').click(function(event) {
		event.preventDefault();
		nodeInstall("COMPUTE");
		return false;
		// for good measure
	});
	
	$('#HA_install_button').click(function(event) {
		event.preventDefault();
		nodeInstall("HA");
		return false;
		// for good measure
	});	

	$("#storage_selection input:radio").on("ifClicked", function() {
		if ($(this).attr("value") == "yes") {
			$("#storage_note").show();
			waiting_nodes_connection("COMPUTE");
		} else {
			$("#storage_note").hide();
		}

	});
	
	$("#ha_selection input:radio").on("ifClicked", function() {
		if ($(this).attr("value") == "yes") {
			$("#ha_note").show();
			waiting_nodes_connectionha("HA");
		} else {
			$("#ha_note").hide();
		}

	});

});

function get_install_text(txt, name) {
	if (txt) {
		$("#" + name + "_status").addClass("text-success");
		$("#" + name + "_install_button").hide();
		$("#" + name + "_uninstallbutton").show();
		$("#" + name + "_statuscolor").removeClass("table-top bg-yellow");
		$("#" + name + "_statuscolor").addClass("table-top bg-green-dark");		
		return "Installed.";
	} else {
		$("#" + name + "_install_button").show();
		$("#" + name + "_uninstallbutton").hide();
		$("#" + name + "_statuscolor").removeClass("table-top bg-green-dark");
		$("#" + name + "_statuscolor").addClass("table-top bg-yellow");
		return "Not Installed.";
	}
}

function installProcess(str) {
	serverID = str.concat("_waiting1");
	successID = str.concat("_success");
	errorID = str.concat("_error");
	buttonID = str.concat("_install_button");

	$('#' + serverID).waiting({
		className : 'waiting-circles',
		elements : 8,
		radius : 20,
		auto : true
	});
	$('#' + buttonID).hide();

	$('#' + serverID).show();
	install_text = str.concat("_install_text");
	$('#' + install_text).show();

	$.ajax({
		type : "GET",
		url : "/servers/install/" + str,
		data : str,
		dataType : 'text',
		async : true,
		crossDomain : "true",
		success : function(response) {
			var res = JSON.parse(response);
			console.log(res.success);
			if (res.success) {
				var ser = JSON.parse(res.data);
				$('#' + serverID).hide();
				$('#' + install_text).hide();
				$('#' + successID).show();
				//	if (i < 2) {
				//	$("." + progress).css('width', '100%');
				//	}
				$("#" + str + "_status").text(get_install_text(ser.Install, ser.Name));
				$("#" + str + "_installdate").text(ser.InstallDate);
				$("#" + str + "_updatedate").text(ser.UpdateDate);
				$("#" + name + "_status").addClass("text-success");
				$("#" + name + "_install_button").hide();
				$("#" + name + "_uninstallbutton").show();
				$("#" + name + "_dash_success").show();
			} else {
				$('#' + serverID).hide();
				$('#' + install_text).hide();
				$('#' + errorID).show();
				$("#" + name + "_dash_error").show();
				if (i == 0) {
					$("." + progress).css('width', '100%');
					$('#change_' + str).removeClass('progress-bar-info').addClass('progress-bar-danger');
				}
			}
		},
		error : function(xhr, status) {
			alert("error " + status);
			$('#' + serverID).hide();
			$('#' + install_text).hide();
			$('#' + errorID).show();
			$("#" + name + "_dash_error").show();		
		}
	});
	return false;
}

function nodeInstall(str) {
    serverID = str.concat("_waiting");
	successID = str.concat("_success");
	errorID = str.concat("_error");
	buttonID = str.concat("_install_button");
	var urlvalue = "";
	
	if(str == "COMPUTE") {	
        str_ip = $("#hostip").val();
        ip = str_ip.split("="); 
        urlvalue = "/nodes/request/" + ip[1];
    } else {
        ip = $("#hosthaip").val();      
     //   ip = str_ip.split("=");
        urlvalue = "/nodes/harequest/" + ip;
    }
    
	$('#' + serverID).waiting({
		className : 'waiting-circles',
		elements : 8,
		radius : 20,
		auto : true
	});
	$('#' + buttonID).hide();
    
	$('#' + serverID).show();
	install_text = str.concat("_install_text");
	$('#' + install_text).show();
    
	$.ajax({
		type : "GET",
	    url : urlvalue,
		data : str,
		dataType : 'text',
		async : true,
		crossDomain : "true",
		success : function(response) {
			var res = JSON.parse(response);	
			if (res.success) {
			//	var ser = JSON.parse(res.data);
				$('#' + serverID).hide();
				$('#' + install_text).hide();
				$('#' + successID).show();
				if (res.data) {	
				  updateNodesList(res.data);
				}							
			} else {
				$('#' + serverID).hide();
				$('#' + install_text).hide();
				$('#' + errorID).show();
				$("#" + name + "_dash_error").show();		
				     $("#myNodesModal").modal('show');
				     $("#res_nodes_msg").text(res.error);
			}
		},
		error : function(xhr, status) {
			alert("error " + status);
			$('#' + serverID).hide();
			$('#' + install_text).hide();
			$('#' + errorID).show();
			$("#" + name + "_dash_error").show();			
		}
	});
	return false;
}


function install_check(i, servers) {
	console.log(i);
	serverID = servers[i].concat("_waiting1");
	successID = servers[i].concat("_success");
	errorID = servers[i].concat("_error");
	buttonID = servers[i].concat("_install_button");
	textID = servers[i].concat("_install_text");

	$.ajax({
		type : "GET",
		url : "/servers/verify/" + servers[i],
		data : servers[i],
		dataType : 'text',
		async : true,
		crossDomain : "true",
		success : function(response) {
			var res = JSON.parse(response);
			if (res.success) {
				$('#' + buttonID).hide();
				$('#' + serverID).hide();
				$('#' + textID).hide();
				$('#' + successID).show();
			} else {
				$('#' + serverID).hide();
				$('#' + textID).hide();
				$('#' + errorID).hide();
				$('#' + buttonID).show();
			}
			if (i < 1) {
				install_check(i + 1, servers);
			}
		},
		error : function(xhr, status) {
			$('#' + serverID).hide();
			$('#' + textID).hide();
			$('#' + errorID).hide();
			$('#' + buttonID).show();
		}
	});
	return false;
}

//waiting the getting ip for host node
function waiting_nodes_connection(nodes) {
	serverNodeID = nodes.concat("_waiting");
	buttonNodeID = nodes.concat("_install_button");
	networkWarningTextID = nodes.concat("_network_warning_text");
	networkSuccessTextID = nodes.concat("_network_success_text");
	$('#' + serverNodeID).waiting({
		className : 'waiting-circles',
		elements : 8,
		radius : 20,
		auto : true
	});
	$.ajax({
		type : "GET",
		url : "/servers/getIP",
		data : nodes,
		dataType : 'text',
		async : true,
		success : function(response) {
			var res = JSON.parse(response);
			console.log(res);
		   
			if (res.ip) {
				$('#' + serverNodeID).hide();
				$('#' + buttonNodeID).show();
				$('#' + networkSuccessTextID).show();
				$('#' + networkWarningTextID).hide();
				$("#hostip").val(res.ipvalue);
			//	$("#hosthaip").val(res.ipvalue);
			}
		},
		error : function(xhr, status) {
			$('#' + serverNodeID).show();
			$('#' + buttonNodeID).hide();
			$('#' + networkSuccessTextID).hide();
			$('#' + networkWarningTextID).show();
		}
	});
	return false;
}

//waiting the getting ip for host node
function waiting_nodes_connectionha(nodes) {
	serverNodeID = nodes.concat("_waiting");
	buttonNodeID = nodes.concat("_install_button");
	networkWarningTextID = nodes.concat("_network_warning_text");
	networkSuccessTextID = nodes.concat("_network_success_text");
	$('#' + serverNodeID).waiting({
		className : 'waiting-circles',
		elements : 8,
		radius : 20,
		auto : true
	});
	$.ajax({
		type : "GET",
		url : "/servers/getIP",
		data : nodes,
		dataType : 'text',
		async : true,
		success : function(response) {
			var res = JSON.parse(response);		
			var r = (res.ipvalue).split("=");
			console.log(r[1]);
			$.ajax({
				type : "GET",
				url : "/servers/getHAOptions/"+ r[1],
				data : r[1],
				dataType : 'text',
				async : true,
				success : function(response) {
					console.log(response);			
					
					var res = JSON.parse(response);
					if (res.success) {
					$('#' + serverNodeID).hide();
					$('#' + buttonNodeID).show();
					$('#' + networkSuccessTextID).show();
					$('#' + networkWarningTextID).hide();
			//		$("#hostip").val(res.ipvalue);
					$("#hosthaip").val(r[1]);
				   } else {
				     $("#myModal").modal('show');
				     $("#res_msg").text(res.error);
				   }				
				},
				error : function(xhr, status) {
			        $('#' + serverNodeID).show();
					$('#' + buttonNodeID).hide();
					$('#' + networkSuccessTextID).hide();
					$('#' + networkWarningTextID).show();
				}
			});			
		},
		error : function(xhr, status) {
			$('#' + serverNodeID).show();
			$('#' + buttonNodeID).hide();
			$('#' + networkSuccessTextID).hide();
			$('#' + networkWarningTextID).show();
		}
	});
	return false;
}

function getHADatas(nodes) {
  $.ajax({
		type : "GET",
		url : "/servers/getHAOptions",
		data : nodes,
		dataType : 'text',
		async : true,
		success : function(response) {
			console.log("-------------------");
			console.log(response);
		},
		error : function(xhr, status) {
			console.log("-------------------");
		}
	});
}


function updateNodesList(res) {
   var obj = JSON.parse(res);
    $("#nodes_list").append(
		'<tr style="height: 50px;">\
		<td class="success">' + obj.IP + '</td>\
		<td class="warning">' + get_install_text(obj.Install) + '</td>\
		<td class="danger">' + obj.InstallDate + '</td>\
		<td class="active">' + obj.UpdateDate + '</td>\
		</tr>');
}


function dashboard(stype) {
   
	$('#' + stype +"_dash_waiting1").waiting({
		className : 'waiting-circles',
		elements : 24,
		radius : 60,
		auto : true
	});
	
   $.ajax({
		type : "GET",
		url : "/dashboard/"+stype,
		dataType : 'text',
		async : true,
		crossDomain : "true",
		success : function(response) {
		        var flag = true;
		        var text;	
		        $('#' + stype +"_dash_waiting1").hide();
		        console.log(response);	       
				var res = JSON.parse(response);
				if(res.success) {
				console.log(res.data);  
				var jsondata = JSON.parse(res.data);				
				 var installkey = jsondata.packages;
		 
					for(p in installkey)
					{
					    $("#coll"+stype).append(
					    '<div class="col-sm-6">\
				    		<h4>' + p + '</h4>\
							<table class="features-table">\
							<thead>\
							<tr>\
								<td></td>\
								<td>Install Status</td>\
								<td>Running Status</td>\
							</tr>\
							</thead>\
							<tbody id="t_'+stype+'_'+p+'">\
							</tbody>\
					      </table>\
						</div>');	
						
					    var in_verify = true;
					    var img;
 						elementArray = installkey[p];
 						for(v in elementArray)
					    {
 					 	  v1 = elementArray[v]; 					
 						  if(v1 == "false") {
 						     in_verify = false;
 						     img = "cross";
 						  } else { 						    
 						     img = "check";
 						  }
 						  $("#t_"+stype+'_'+p).append('\
 						        <tr>\
									<td>' + v + '</td>\
									<td><img src="/static_source/images/'+ img +'.png" width="16" height="16" alt="'+img+'"></td>\
									<td id="running_status_'+stype+'_'+v+'"></td>\
								</tr>');
					   }
					   console.log(in_verify);
						if(in_verify == true) {
 						   text = "Installed";
 						   badge = "timeline-badge success";
 						 } else {
 						   text = "Not Installed";
 						   badge = "timeline-badge danger";
 						 }     	
 						 
 					if(flag==true) { 						
 						$("#d_"+stype).append(				 
							'<li>\
								<div class="' + badge + '">\
									<i class="fa fa-check"></i>\
								</div>\
								<div class="timeline-panel">\
									<div class="timeline-heading">\
										<h4 class="timeline-title">' + p + '</h4>'+										
									'</div>\
									<div class="timeline-body">\
										<p>'+ text + '</p>\
									</div>\
									<div class="timeline-body">\
										<p id="' + stype + '_' + p + '_status"></p>\
									</div>\
								</div>\
							</li>'); 
							flag = false;
					   } else {
					      $("#d_"+stype).append(				 
							'<li class="timeline-inverted">\
								<div class="' + badge + '">\
									<i class="fa fa-check"></i>\
								</div>\
								<div class="timeline-panel">\
									<div class="timeline-heading">\
										<h4 class="timeline-title">' + p + '</h4>'+										
									'</div>\
									<div class="timeline-body">\
									   <p>'+ text + '</p>\
									</div>\
									<div class="timeline-body">\
										<p id="' + stype + '_' + p + '_status"></p>\
									</div>\
								</div>\
							</li>'); 
							flag = true;
					   } 	
					}
				 setStatus(jsondata, stype);
			   }        			
				},
				error : function(xhr, status) {
					console.log(status);
			}
	});
	
	}
	

function setStatus(json, stype) {
    var servicekey = json.services;
	var status;			 
	for(p in servicekey)
	{
		var in_verify = true;
		var img;
 		elementArray = servicekey[p];
 		for(v in elementArray)
		{	
 		  v1 = elementArray[v]; 					
 		   if(v1 == "false") {
 		     in_verify = false;
 			 img = "cross";
 		  } else {
 			 img = "check";
 		  }
 		  if(getService(v) == true) {
 		  	if(v == "snowflake") {
 		    	v = "megamsnowflake";
 		  	}  
 		  	$("#running_status_"+stype+"_"+v).html('<img src="/static_source/images/'+ img +'.png" width="16" height="16" alt="'+img+'">');
		 }	else {
		   $("#running_status_"+stype+"_"+v).html('No service');
		 }
		if(in_verify == true) {
 		   status = "Running";
 		   badge = "timeline-badge success";
 		} else {
 		   status = "Not Running";
 		   badge = "timeline-badge danger";
 		}     						 
 		$("#"+stype+"_"+p+"_status").text(status);	
 		
	} 
 }	
}

function getService(serviceName) {
 var value;
 console.log(serviceName);
 switch (serviceName) {
    case "megamcommon":
        value = false;
        break;
    case "ruby2.0":
        value = false;
        break;
    case "debmirror":
        value = false;
        break;
   default:
        value = true;
        break;
  }
  return value;
}
