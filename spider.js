var casper = require('casper').create({
    verbose: true,
    logLevel: "error"
});
var fs = require('fs');

casper.userAgent('Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36');
casper.start('https://www.peachjohn.co.jp/pjitem/list/?_al=16wi_sale&outlet=only&count=2000&_dm=2&itemSrcFlt=1&sort=6&page=1');

// waitForSelector
casper.wait(150000, function () {
    var images = casper.evaluate(function () {
        var pics = $('.pic img');

        return Array.prototype.map.call(pics, function (pic) {
            return {
                src: $(pic).attr('src'),
            };
        });
    });

    var result = '';
    for (var index = 0; index < images.length; index++) {
        result += images[index].src.split('?')[0] + '?wid=640&op_usm=1%2C1%2C10%2C0&resMode=sharp2\r\n';
    }

    fs.write(fs.pathJoin(fs.workingDirectory, 'output', 'meta'), result, 'w');
});


/*
casper.waitForSelector('.pic img', function () {
    var images = casper.evaluate(function () {
        var pics = $('.pic img');

        return Array.prototype.map.call(pics, function (pic) {
            return {
                src: $(pic).attr('src').split('?')[0] + '?wid=640&op_usm=1%2C1%2C10%2C0&resMode=sharp2',
            };
        });
    });

    var index = 1;
    casper.eachThen(images, function (response) {
        var url = response.data.src;

        casper.open(url).then(function () {
            casper.captureSelector('./output/' + index++ + '.png', 'img');
        });
    });
});
*/

casper.run();