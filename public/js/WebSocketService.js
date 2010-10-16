var WebSocketService = function(model, webSocket) {
	var webSocketService = this;
	
	var webSocket = webSocket;
	var model = model;
	
	this.hasConnection = false;
	
	this.welcomeHandler = function(data) {
		webSocketService.hasConnection = true;
		
		model.userTadpole.id = data.Id;
		model.tadpoles[data.Id] = model.tadpoles[-1];
		delete model.tadpoles[-1];
		
		$('#chat').initChat();
	};
	
	this.updateHandler = function(data) {
		var newtp = false;

		if(!model.tadpoles[data.Id]) {
			newtp = true;
			model.tadpoles[data.Id] = new Tadpole();
			model.arrows[data.Id] = new Arrow(model.tadpoles[data.Id], model.camera);
		}
		
		var tadpole = model.tadpoles[data.Id];
		
		if(data.Id == model.userTadpole.id) {
			if(!model.userTadpole.Name) {
				tadpole.name = data.Name;
			}
			return;
		} else {
			tadpole.name = data.Name;
		}
		
		if(newtp) {
			tadpole.x = data.X;
			tadpole.y = data.Y;
		} else {
			tadpole.targetX = data.X;
			tadpole.targetY = data.Y;
		}
		
		tadpole.angle = data.Angle;
		tadpole.momentum = data.Momentum;
		
		tadpole.timeSinceLastServerUpdate = 0;
	}

	this.contentHandler = function(data) {
		model.content[data.Id] = new Content(data);
	}
	
	this.messageHandler = function(data) {
		var tadpole = model.tadpoles[data.Id];
		if(!tadpole) {
			return;
		}
		tadpole.timeSinceLastServerUpdate = 0;
		tadpole.messages.push(new Message(data.Message));
	}
	
	this.closedHandler = function(data) {
		if(model.tadpoles[data.Id]) {
			delete model.tadpoles[data.Id];
			delete model.arrows[data.Id];
		}
	}
	
	this.processMessage = function(data) {
		var fn = webSocketService[data.type.toLowerCase() + 'Handler'];
		if (fn) {
			fn(data.data);
		}
	}
	
	this.connectionClosed = function() {
		webSocketService.hasConnection = false;
		$('#cant-connect').fadeIn(300);
	};
	
	this.sendUpdate = function(tadpole) {
		var sendObj = {
			Update: {
				X: tadpole.x,
				Y: tadpole.y,
				Angle: tadpole.angle,
				Momentum: tadpole.momentum
			}
		};
		
		if(tadpole.name) {
			sendObj.Update['Name'] = tadpole.name;
		}
		
		webSocket.send(JSON.stringify(sendObj));
	}
	
	this.sendMessage = function(msg) {
		var regexp = /name: ?(.+)/i;
		if(regexp.test(msg)) {
			model.userTadpole.name = msg.match(regexp)[1];
			return;
		}
		
		var sendObj = {
			Message: { Message: msg }
		};
		
		webSocket.send(JSON.stringify(sendObj));
	}
}
