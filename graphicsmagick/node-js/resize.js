// Init GraphicsMagick and load image
var gm = require('gm'),
	start = +new Date(),
	image = gm('test.jpg');

(function execute(i) {
	// Save thumbnail
	image.thumb(200, 100, 'test-thumb-' + i + '.jpg', 79, function() {
		if (i < 19) {
			execute(++i);
		} else {
			console.log('Time taken: ', ((+new Date() - start) / 1000) + 's');
		}
	});	
})(0);