var Content = function(data) {
	var content = this;

	this.id = data.Id;
	this.x = data.X;
	this.y = data.Y;
	this.size = 0;
	this.targetSize = 10;
	this.color = '200,200,100';

	function tween(value, target, rate) {
		var r = value;
		if (value < target) {
			r += target*rate - value*rate;
			if (r > target * 0.95)
				r = target;
		}
		return r
	}

	this.draw = function(context) {
		// animate
		content.size = tween(content.size, content.targetSize, 0.1);

		var opacity = 1.0;
		context.fillStyle     = 'rgba('+content.color+','+opacity+')';
		context.shadowColor   = 'rgba('+content.color+','+opacity*0.7+')';
		context.shadowOffsetX = 0;
		context.shadowOffsetY = 0;
		context.shadowBlur    = 10;

		// Draw circle
		context.beginPath();
		context.arc(content.x, content.y, content.size, 0, Math.PI*2, true);
		context.closePath();
		context.fill();
	};
}
