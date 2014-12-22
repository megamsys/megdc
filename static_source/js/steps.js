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
		nodeInstall();
		return false;
		// for good measure
	});

	$("#ha_selection input:radio").on("ifClicked", function() {
		if ($(this).attr("value") == "yes") {
			$("#ha_note").show();
		} else {
			$("#ha_note").hide();
		}
		$(this).attr('disabled', true);

	});

	$("#storage_selection input:radio").on("ifClicked", function() {
		if ($(this).attr("value") == "yes") {
			$("#storage_note").show();
			waiting_nodes_connection("COMPUTE");
		} else {
			$("#storage_note").hide();
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
			//if (i == 0) {
			//	$("." + progress).css('width', '100%');
			//	$('#change_' + str).removeClass('progress-bar-info').addClass('progress-bar-danger');
			//}
		}
	});
	return false;
}

function nodeInstall() {
	serverID = "COMPUTE_waiting1";
	successID = "COMPUTE_success";
	errorID = "COMPUTE_error";
	buttonID = "COMPUTE_install_button";
    str = $("#hostip").val();
    
	$('#' + serverID).waiting({
		className : 'waiting-circles',
		elements : 8,
		radius : 20,
		auto : true
	});
	$('#' + buttonID).hide();
    
	$('#' + serverID).show();
	install_text = "COMPUTE_install_text";
	$('#' + install_text).show();

	$.ajax({
		type : "GET",
		url : "/nodes/request/" + str,
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
			} else {
				$('#' + serverID).hide();
				$('#' + install_text).hide();
				$('#' + errorID).show();
				$("#" + name + "_dash_error").show();				
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

/*function installProcess(i) {
 var servers = [ "MEGAM", "COBBLER", "OPENNEBULA", "OPENNEBULAHOST", "STORAGE" ];
 serverID = servers[i].concat("_waiting1");
 successID = servers[i].concat("_success");
 errorID = servers[i].concat("_error");
 buttonID = servers[i].concat("_install_button");

 if (i > 1) {
 console.log("opennebula entry");
 $('#' + serverID).waiting({
 className : 'waiting-circles',
 elements : 8,
 radius : 20,
 auto : true
 });
 $('#' + buttonID).hide();
 } else {
 $('#' + serverID).waiting({
 className : 'waiting-circles',
 elements : 8,
 radius : 20,
 auto : true
 });
 // $('#' + serverID).waiting({
 // className : 'waiting-blocks',
 // elements : 5,
 // speed : 200,
 // auto : true
 // });
 progress = servers[i].concat("_PROGRESS");
 $("#" + progress).show();
 $("." + progress).css('width', '50%');
 log = servers[i].concat("_LOG");
 $("#" + log).show();

 }

 $('#' + serverID).show();
 install_text = servers[i].concat("_install_text");
 $('#' + install_text).show();

 $.ajax({
 type : "GET",
 url : "/servers/" + servers[i],
 data : servers[i],
 dataType : 'text',
 async : true,
 crossDomain : "true",
 success : function(response) {
 var res = JSON.parse(response);
 console.log(res.success);
 if (res.success) {
 $('#' + serverID).hide();
 $('#' + install_text).hide();
 $('#' + successID).show();
 if (i < 2) {
 $("." + progress).css('width', '100%');
 }
 if (i == 0) {
 $("#wzdButtons").show();
 //installProcess(i + 1);
 }
 //if (i == 1) {
 //$("#wzdButtons").show();
 //}
 } else {
 $('#' + serverID).hide();
 $('#' + install_text).hide();
 $('#' + errorID).show();
 if (i == 0) {
 $("." + progress).css('width', '100%');
 $('#change_' + servers[i]).removeClass('progress-bar-info')
 .addClass('progress-bar-danger');
 }
 }
 },
 error : function(xhr, status) {
 alert("error " + status);
 $('#' + serverID).hide();
 $('#' + install_text).hide();
 $('#' + errorID).show();
 if (i == 0) {
 $("." + progress).css('width', '100%');
 $('#change_' + servers[i]).removeClass('progress-bar-info')
 .addClass('progress-bar-danger');
 }
 }
 });
 return false;
 }

 function install_check(pname) {
 serverID = pname.concat("_waiting1");
 successID = pname.concat("_success");
 errorID = pname.concat("_error");
 buttonID = pname.concat("_install_button");
 textID = pname.concat("_install_text");

 $.ajax({
 type : "GET",
 url : "/servers/verify/" + pname,
 data : pname,
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
 */

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
			if (res.ip) {
				$('#' + serverNodeID).hide();
				$('#' + buttonNodeID).show();
				$('#' + networkSuccessTextID).show();
				$('#' + networkWarningTextID).hide();
				$("#hostip").val(res.ipvalue);
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
