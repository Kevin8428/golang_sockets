
var initHost = function() {

  var conn,
  questions = document.getElementById('host-questions');

  function appendQuestion(question) {
    questions.appendChild(question)
  }

  if (window["WebSocket"]) {
    conn = new WebSocket("ws://localhost:8081/admin_ws");
    conn.onclose = function(evt) {
      var comment = document.createElement('div');
      comment.innerHTML = '<b>Connection closed.<\/b>';
      appendQuestion(comment);
    }
    conn.onmessage = function(evt) {
      var comment = document.createElement('div');
      comment.innerHTML = evt.data;
      appendQuestion(comment);
    }
  } else {
    var comment = document.createElement('div');
    comment.innerHTML = '<b>Your browser does not support WebSockets.<\/b>';
    appendQuestion(comment);
  }
};

var initAdmin = function() {
  var conn,
  msg = document.getElementById('submit-message'),
  log = document.getElementById('all-messages');

  function appendToChat(msg) {
    log.appendChild(msg)
  }

  document.getElementById('submit-message-form').addEventListener('submit', function(e){
    e.preventDefault();
      if (!conn) {
        console.log("no connection")
        return false;
      }
      if (!msg.value) {
        console.log("no value")
        return false;
      }
      conn.send(msg.value);
      msg.value = "";
      return false
  });

  if (window["WebSocket"]) {
    conn = new WebSocket("ws://localhost:8081/customer_ws");
    conn.onclose = function(evt) {
      var comment = document.createElement('div');
      comment.innerHTML = '<b>Connection closed.<\/b>';
      appendToChat(comment)
    }
    conn.onmessage = function(evt) {
      var comment = document.createElement('div');
      comment.innerHTML = evt.data;
      appendToChat(comment)
    }
  } else {
    var comment = document.createElement('div');
    comment.innerHTML = '<b>Your browser does not support WebSockets.<\/b>';
    appendToChat(comment)
  }
  //////////////////////////////////////////////////////////////////////////////////
  ///////////////////////////// ADMIN WEBSOCKET ////////////////////////////////////
  //////////////////////////////////////////////////////////////////////////////////
  var adminConn,
  hostMessage = document.getElementById('host-questions'),
  adminLog = document.getElementById('admin-log');

  function appendToAdmin(hostMessage) {
    adminLog.appendChild(hostMessage)
  }

  document.getElementById('show-host-form').addEventListener('submit', function(e){
    e.preventDefault();
  if (!adminConn) {
    console.log("no connection");
    return false;
  }
  if (!hostMessage.value) {
    console.log("no hostMessage");
    return false;
  }
  adminConn.send(hostMessage.value);
  hostMessage.value = '';
  console.log('submitting')
  return false
  });

  // create chat connection and render messages to page
  if (window["WebSocket"]) {
      adminConn = new WebSocket("ws://localhost:8081/admin_ws");
      adminConn.onopen = function() {
        console.log('admin websocket connection made!')
      }
      adminConn.onclose = function(evt) {
        var comment = document.createElement('div');
        comment.innerHTML = '<b>Admin connection closed.<\/b>';
          appendToAdmin(comment)
      }
      adminConn.onmessage = function(evt) {
        console.log("received message");
        var comment = document.createElement('div');
        comment.innerHTML = evt.data;
        appendToAdmin(comment)
      }
  } else {
    var comment = document.createElement('div');
    comment.innerHTML = '<b>browser does not support websockets.<\/b>';
    appendToAdmin(comment)
  }  
}

var initCustomerPage = function() {
  var conn,
  msg = document.getElementById('submit-message'),
  log = document.getElementById('all-messages');

  function appendToChat(msg) {
    log.appendChild(msg)
  }

  document.getElementById('submit-message-form').addEventListener('submit', function(e){
    e.preventDefault();
      if (!conn) {
        console.log("no connection")
        return false;
      }
      if (!msg.value) {
        console.log("no value")
        return false;
      }
      conn.send(msg.value);
      msg.value = "";
      return false
  });

  if (window["WebSocket"]) {
    conn = new WebSocket("ws://localhost:8081/customer_ws");
    conn.onclose = function(evt) {
      var comment = document.createElement('div');
      comment.innerHTML = '<b>Connection closed.<\/b>';
      appendToChat(comment)
    }
    conn.onmessage = function(evt) {
      var comment = document.createElement('div');
      comment.innerHTML = evt.data;
      appendToChat(comment)
    }
  } else {
    var comment = document.createElement('div');
    comment.innerHTML = '<b>Your browser does not support WebSockets.<\/b>';
    appendToChat(comment)
  }
}

window.onload = function(){
  var isAdminPage = document.getElementById('admin-log') != null
  var isHostPage = document.getElementById('host-questions') != null

  if (isAdminPage) {
    initAdmin();
  } else if (isHostPage) {
    initHost();
  } else {
    initCustomerPage()
  }
}