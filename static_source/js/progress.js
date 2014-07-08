/*
 ** Copyright [2012-2013] [Megam Systems]
 **
 ** Licensed under the Apache License, Version 2.0 (the "License");
 ** you may not use this file except in compliance with the License.
 ** You may obtain a copy of the License at
 **
 ** http://www.apache.org/licenses/LICENSE-2.0
 **
 ** Unless required by applicable law or agreed to in writing, software
 ** distributed under the License is distributed on an "AS IS" BASIS,
 ** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 ** See the License for the specific language governing permissions and
 ** limitations under the License.
 */

//$(document).ready(function(){
//  var progresspump = setInterval(function(){
/* query the completion percentage from the server */
//$.get("/servers/1", function(data){
//  $..load("/servers/1", function(response, status, xhr) {	
/* update the progress bar width */
//   $(".progress-bar").css('width',responsedata+'%');
/* and display the numeric value */
//   $(".progress-bar").html(data+'%');
/* test to see if the job has completed */
//  if(data > 99.999) {
//    clearInterval(progresspump);
// $("#progressouter").removeClass("active");
//     $(".progress-bar").html("Done");
//   }
//  })
// }, 1000);});

var socket;

$(document).ready(
		function() {
			// Create a socket
			socket = new WebSocket('ws://' + window.location.host
					+ '/servers/join?uname=megam');
			// Message received on the socket
			socket.onmessage = function(event) {
				var data = $.parseJSON(event.data);
				console.log(data);
				console.log(parseInt(data.Content));
				$(".progress-bar").css('width', parseInt(data.Content) + '%');
			};

			// Send messages.
			var postConecnt = function() {
				// var uname = $('#uname').text();
				// var content = $('#sendbox').val();
				// socket.send("megam");
				// $('#sendbox').val("");
			}

			// $('#sendbtn').click(function () {
			postConecnt();
			// });
		});