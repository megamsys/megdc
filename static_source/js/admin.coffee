$(document).ready ->

  # chat demo
  makeMessageString = (content)->
    "<li class='right'><img src='placeholders/avatars/9.jpg' class='img-circle'><div class='message'>#{content}</div></li>"
  makeReplyString = (content)->
    "<li><img src='placeholders/avatars/avatar.jpg' class='img-circle' width='26'><div class='message'>#{content}</div></li>"

  addMessage = (content, reply=false)->
    messages = $('.chat-messages')
    message = 
      if reply 
      then makeReplyString(content) 
      else makeMessageString(content)
    messages.append(message)
    messages.scrollTop( messages.height() )

  if $('.messenger').length > 0
    addMessage("Hey, how are you?")
    setTimeout ->
      addMessage("You like Sugoi Admin?")
    , 5000
    setTimeout ->
      addMessage("Quiet slick, isn't it? Have a look at widgets page and remember - there's a lot more coming soon!")
    , 12000

    $("#message-input").keyup (e)->
      if e.keyCode == 13
        addMessage $(@).val(), true
        $(@).val("")

    $("#chat-toggle").on 'click', (e)->
      $('.messenger-body').toggleClass('open');
      $('#chat-toggle .glyphicon').toggleClass("glyphicon-chevron-down glyphicon-chevron-up")


  # notifications
  # visit: http://notifyjs.com/ 
  if $('#inbox-page').length > 0
    $.notify "You've got 4 new messages", 'success',
      autoHide: true
      autoHideDelay: 5000
      arrowShow: false

  if $('#dashboard-page').length > 0
    $.notify "Welcome back John", 'info',
      autoHide: true
      autoHideDelay: 5000
      arrowShow: false
  


  # Gallery image alignment setup 
  # visit: http://sapegin.github.io/jquery.mosaicflow/
  $(".mosaicflow__item").each ->
    s = Snap @.querySelector('svg')
    path = s.select 'path'

    pathConfig = 
      from: path.attr 'd'
      to: @.getAttribute 'data-path-hover'

    @.addEventListener 'mouseenter', ->
      path.animate { 'path' : pathConfig.to }, 300, mina.easeinout

    @.addEventListener 'mouseleave', ->
      path.animate { 'path' : pathConfig.from }, 300, mina.easeinout


  # gallery image upload setup
  # visit: http://www.dropzonejs.com/
  if $('.dropzone').length > 0
    Dropzone.options.demoUpload =
      paramName: "file"
      maxFilesize: 2
      addRemoveLinks: true



  # sidebar chart with easypiecharts
  # visit: https://github.com/rendro/easy-pie-chart
  sidebarChart = $('.sidebar-chart')
  sidebarChart.easyPieChart {
    barColor: "#2CC0D5"
    trackColor: "rgba(255,255,255,.06)"
    lineWidth: 10
    animate: 600
    lineCap: "square"
    size: 140
    onStart: (from, to)->
      percentage = $(@.el).find('.percentage')
      $({val: from}).animate {val: to}, {
          duration: 600
          easing: 'swing'
          step: ()->
            percentage.text( Math.round(@.val) + "%" )
        }
  }

  # animate sidebar chart every 5 seconds with random value
  animateSidebarChart = ()->
    random = Math.round(Math.random() * (100 - 1) + 1) # get random percentage
    setTimeout ()->
      $('.sidebar-chart').data('easyPieChart').update(random)
      animateSidebarChart()
    , 5000
  animateSidebarChart()


  # make messages selectable (aka prevent click to go to msg)
  inboxMessages = $('.messages-selectable')
  if inboxMessages.length > 0
    inboxMessages.on 'click', 'a', (e)->
      e.preventDefault()
      $(@).toggleClass 'active'




  # List of tasks
  # visit: http://farhadi.ir/projects/html5sortable/
  listSortable = $('.list-editable')
  if listSortable.length > 0
    listItem = $('#item-template').clone()
    taskInput = $('#task-content')

    # toggle input field for task
    # and focus on the input
    $('#task-toggle').on 'click', (e)->
      e.preventDefault()
      taskInput.toggleClass('opacity-1')
      if taskInput.hasClass('opacity-1')
        taskInput.focus()

    # add task to DOM/save to db db
    taskInput.on 'keyup', (e)->
      if e.keyCode == 13
        e.preventDefault()
        value = $('#task-content').val()
        task = listItem.clone()
        $('.body', task).text(value)
        task.removeAttr('style')
        task.removeAttr('id')
        listSortable.prepend(task)
        $('#task-content').val('')

        # add checkbox styles
        $('input', task).iCheck {
          checkboxClass: "icheckbox_flat-blue"
        }

    listSortable.on 'click', '.remove-box', (e)->
      e.preventDefault()
      $(this).parent().remove()



  # iCHECK - better checkboxes
  # visit: https://github.com/fronteed/iCheck
  if $('input.icheck').length > 0
    $('input.icheck').iCheck {
      checkboxClass: "icheckbox_flat-blue"
      radioClass: 'iradio_flat-blue'
    }



  # Better select boxes
  # visit: https://github.com/silviomoreto/bootstrap-select
  if $(".selectpicker").length > 0
    $('.selectpicker').selectpicker()



  # SOCIAL FEED
  # visit: https://github.com/christianv/jquery-lifestream/
  if $("#social-feed").length > 0
    $("#social-feed").lifestream {
      limit: 30
      list:[
        {
          service: "twitter",
          user: "twbootstrap"
        },
        {
          service: "stackoverflow",
          user: '117193'
        },
        {
          service: "rss",
          user: "http://qz.com/feed/"
        },
        {
          service: "blogger",
          user: "googleblog"
        },
        {
          service: "dribbble",
          user: 'dribbble'
        }
      ]
    }

    # add social icons
    $('#social-feed li.livestream-twitter').each ->
      twitterIcon = '<i class="fc-social-twitter"></i>'
      $(this).prepend(twitterIcon)



  # Events calendar
  # visit: http://www.vissit.com/jquery-event-calendar-plugin-english-version
  if $('#events-calendar').length > 0
    today    = moment().utc().format("YYYY-MM-DD HH:MM:SS")
    tomorrow = moment().utc().add('days', 1).format("YYYY-MM-DD HH:MM:SS")
    someday  = moment().utc().add('days', 3).format("YYYY-MM-DD HH:MM:SS")

    events = [{
        "date": today
        "type": "meeting"
        "title": "Mattis Cras"
        "description": "Ipsum Consectetur Etiam Nibh"
        "url": "http://www.event1.com/"
      },
      {
        "date": tomorrow
        "type": "meeting"
        "title": "Pellentesque Parturient Dolor"
        "description": "Donec ullamcorper nulla non metus auctor fringilla."
        "url": "http://www.event1.com/"
      },
      {
        "date": someday
        "type": "meeting"
        "title": "Ligula"
        "description": "Nulla vitae elit libero, a pharetra augue."
        "url": "http://www.event1.com/"
    }]

    $("#events-calendar").eventCalendar {
      jsonData: events
      jsonDateFormat: 'human'
    }


  # Morris charts
  # visit: https://github.com/oesmith/morris.js/
  if $('a[data-toggle="tab"][href="#statistics"]').length > 0
    stats_shown = false # prevent duplicating
    $('a[data-toggle="tab"][href="#statistics"]').on 'shown.bs.tab', (e)->
      if not stats_shown
        stats_shown = true

        visitorsChart = Morris.Line {
          element: 'visitors-chart'
          pointFlagColors: ["#2CC0D5"]
          lineColors: ["#2CC0D5"]
          resize: true
          data: [
            { date: moment().utc().format("YYYY-MM-DD"), value: 204 }
            { date: moment().utc().add('days', 1).format("YYYY-MM-DD"), value: 155 }
            { date: moment().utc().add('days', 2).format("YYYY-MM-DD"), value: 220 }
            { date: moment().utc().add('days', 3).format("YYYY-MM-DD"), value: 201 }
            { date: moment().utc().add('days', 4).format("YYYY-MM-DD"), value: 198 }
            { date: moment().utc().add('days', 5).format("YYYY-MM-DD"), value: 287 }
            { date: moment().utc().add('days', 6).format("YYYY-MM-DD"), value: 192 }
          ]
          xkey: 'date'
          ykeys: ['value']
          labels: ['Visitors']
          xLabelFormat: (x)->
            return moment(x).format("DD MMM")
          dateFormat: (x)->
            return moment(x).format("dddd DD MMM")
        }

        sourceChart = Morris.Bar {
          element: 'source-chart'
          barColors: ["#2CC0D5"]
          resize: true
          data: [
            { y: 'Search engines', a: 36 }
            { y: 'Social network', a: 29 }
            { y: 'Ad campaign', a: 24 }
            { y: 'Direct traffic', a: 10 }
            { y: 'Other', a: 11 }
          ]
          xkey: 'y'
          ykeys: ['a']
          xLabelAngle: 10
          labels: ['Source']
        }

        deviceChart = Morris.Donut {
          element: 'device-chart'
          colors: ["#2CC0D5", "#37D3EA", "#3BE2FB", "#81F6FF", "#A9F8FF"]
          resize: true
          data: [
            {label: 'iPhone', value: 36}
            {label: 'iPhone 3G', value: 29}
            {label: 'iPhone 3GS', value: 24}
            {label: 'iPhone 4', value: 10}
            {label: 'iPhone 5', value: 11}
          ]
          formatter: (y) -> return y + "%"
        }

        # Sparklines
        # visit: http://omnipotent.net/jquery.sparkline/#s-about 
        twitter = [12, 15, 8, 10, 11, 9, 10]
        $('.twitter-sparkline').sparkline twitter, {
          type: 'bar'
          barColor: '#00ACED'
          height: 22
          barWidth: 8
        }

        facebook = [3, 4, 1, 4, 0, 2, 8]
        $('.facebook-sparkline').sparkline facebook, {
          type: 'bar'
          barColor: '#3b5998'
          height: 22
          barWidth: 8
        }

        google = [3, 8, 10, 2, 6, 2, 8]
        $('.google-sparkline').sparkline google, {
          type: 'bar'
          barColor: '#DA453D'
          height: 22
          barWidth: 8
        }

  # ui page
  if $(".widgets").length > 0
    data = [12, 15, 8, 10, 11, 9, 10]
    $('.ui-sparkline').sparkline data, {
      type: 'bar'
      barColor: '#00ACED'
      height: 22
      barWidth: 8
    }

    $('.ui-sparkline2').sparkline data, {
      type: 'line'
      spotColor: '#FF1B19'
      defaultPixelsPerValue: 8
      minSpotColor: '#FF1B19'
      maxSpotColor: '#FF1B19'
    }

    $('.ui-sparkline3').sparkline [2,1,-1,0,1,3,6,-5], {
      type: 'tristate'
      posBarColor: '#00ACED'
      negBarColor: '#FF1B19'
      barWidth: 8
    }

    $('.ui-sparkline4').sparkline data, {
      type: 'box'
      barColor: '#00ACED'
      height: 22
      barWidth: 8
    }

    $('.ui-sparkline5').sparkline data, {
      type: 'bullet'
      barColor: '#00ACED'
      height: 22
      barWidth: 8
    }

    $('.ui-sparkline6').sparkline data, {
      type: 'pie'
    }

    uiChart1 = Morris.Line {
      element: 'ui-chart-1'
      pointFlagColors: ["#2CC0D5"]
      lineColors: ["#2CC0D5"]
      resize: true
      data: [
        { date: moment().utc().format("YYYY-MM-DD"), value: 204 }
        { date: moment().utc().add('days', 1).format("YYYY-MM-DD"), value: 155 }
        { date: moment().utc().add('days', 2).format("YYYY-MM-DD"), value: 220 }
        { date: moment().utc().add('days', 3).format("YYYY-MM-DD"), value: 201 }
        { date: moment().utc().add('days', 4).format("YYYY-MM-DD"), value: 198 }
        { date: moment().utc().add('days', 5).format("YYYY-MM-DD"), value: 287 }
        { date: moment().utc().add('days', 6).format("YYYY-MM-DD"), value: 192 }
      ]
      xkey: 'date'
      ykeys: ['value']
      labels: ['Visitors']
      xLabelFormat: (x)->
        return moment(x).format("DD MMM")
      dateFormat: (x)->
        return moment(x).format("dddd DD MMM")
    }

    uiChart2 = Morris.Bar {
      element: 'ui-chart-2'
      barColors: ["#2CC0D5"]
      resize: true
      data: [
        { y: 'Search engines', a: 36 }
        { y: 'Social network', a: 29 }
        { y: 'Ad campaign', a: 24 }
        { y: 'Direct traffic', a: 10 }
        { y: 'Other', a: 11 }
      ]
      xkey: 'y'
      ykeys: ['a']
      xLabelAngle: 10
      labels: ['Source']
    }

    uiChart3 = Morris.Donut {
      element: 'ui-chart-3'
      colors: ["#2CC0D5", "#37D3EA", "#3BE2FB", "#81F6FF", "#A9F8FF"]
      resize: true
      data: [
        {label: 'iPhone', value: 36}
        {label: 'iPhone 3G', value: 29}
        {label: 'iPhone 3GS', value: 24}
        {label: 'iPhone 4', value: 10}
        {label: 'iPhone 5', value: 11}
      ]
      formatter: (y) -> return y + "%"
    }

    uiPieChart = $('.ui-piechart')
    uiPieChart.easyPieChart {
      barColor: "#2CC0D5"
      trackColor: "rgba(0,0,0,.06)"
      lineWidth: 10
      animate: 600
      lineCap: "square"
      size: 140
    }

    uiPieChart2 = $('.ui-piechart2')
    uiPieChart2.easyPieChart {
      barColor: "#2CC0D5"
      trackColor: "rgba(0,0,0,.06)"
      lineWidth: 10
      animate: 600
      lineCap: "square"
      size: 140
    }


  # Summernote (WYSIWYG editor)
  summernote = $("#summernote")
  if summernote.length
    summernote.summernote()


  # Maphael World Map
  map = $(".map-container")
  if map.length
    map.mapael
      map:
        name: "world_countries"
        zoom:
          enabled: true,
          maxLebel: 10


  # Hide preloader and show the content
  $('.main-content').addClass('active');
  $('.preloader').fadeOut().remove();

  