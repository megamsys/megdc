$(document).ready(function() {
	var i = 0;
	installProcess(i);

	$('#OPENNEBULA_install_button').click(function(event) {
		event.preventDefault();
		installProcess(2);
		return false; // for good measure
	});
	
	$('#OPENNEBULAHOST_install_button').click(function(event) {
		event.preventDefault();
		installProcess(3);
		return false; // for good measure
	});

    $('#STORAGE_install_button').click(function(event) {
		event.preventDefault();
		installProcess(4);
		return false; // for good measure
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
		} else {
			$("#storage_note").hide();
		}

	});

});

function installProcess(i) {
	var servers = [ "MEGAM", "COBBLER", "OPENNEBULA", "OPENNEBULAHOST", "STORAGE" ];
	serverID = servers[i].concat("_waiting1")
	successID = servers[i].concat("_success")
	errorID = servers[i].concat("_error")
    buttonID = servers[i].concat("_install_button")
	
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
		// cache : true,
		// jsonp : "onJSONPLoad",
		// jsonpCallback: "newarticlescallback",
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
					installProcess(i + 1);
				}
				if (i == 1) {
					$("#wzdButtons").show();
				}
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

function opennebula_install_check() {
	$.ajax({
		type : "GET",
		url : "/servers/verify/OPENNEBULA",
		data : "OPENNEBULA",
		dataType : 'text',
		async : true,
		// cache : true,
		// jsonp : "onJSONPLoad",
		// jsonpCallback: "newarticlescallback",
		crossDomain : "true",
		success : function(response) {
			var res = JSON.parse(response);
			if (res.success) {
				$("#OPENNEBULA_install_button").hide();
				$('#OPENNEBULA_waiting1').hide();
				$('#OPENNEBULA_install_text').hide();
				$('#OPENNEBULA_success').show();
			} else {
				$('#OPENNEBULA_waiting1').hide();
				$('#OPENNEBULA_install_text').hide();
				$('#OPENNEBULA_error').hide();
				$("#OPENNEBULA_install_button").show();
			}
		},
		error : function(xhr, status) {
			$('#OPENNEBULA_waiting1').hide();
			$('#OPENNEBULA_install_text').hide();
			$('#OPENNEBULA_error').hide();
			$("#OPENNEBULA_install_button").show();
		}
	});
	return false;
}

function waiting_nodes_connection(nodes) {
	serverNodeID = nodes.concat("_waiting")
	buttonNodeID = nodes.concat("_install_button")
	networkWarningTextID = nodes.concat("_network_warning_text")
	networkSuccessTextID = nodes.concat("_network_success_text")
	$('#' + serverNodeID).waiting({
		className : 'waiting-circles',
		elements : 8,
		radius : 20,
		auto : true
	});
	$.ajax({
		type : "GET",
		url : "/servers/nodes/"+nodes,
		data : nodes,
		dataType : 'text',
		async : true,
		success : function(response) {			
			var res = JSON.parse(response);
			if(res.ip) {
				$('#' + serverNodeID).hide();				
				$('#' + buttonNodeID).show();
				$('#' + networkSuccessTextID).show();
				$('#' + networkWarningTextID).hide();
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
