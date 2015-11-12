var gulp = require('gulp');
var gutil = require('gulp-util');
var source = require('vinyl-source-stream');
var babelify = require('babelify');
var watchify = require('watchify');
var exorcist = require('exorcist');
var browserify = require('browserify');

// Input file.
watchify.args.debug = true;
var bundler = watchify(browserify('./client/app.js', watchify.args));

// Babel transform
bundler.transform(babelify.configure({
  //sourceMapRelative: 'public/js',
  presets: ["es2015", "stage-0", "react"]
}));

// On updates recompile
bundler.on('update', bundle);

function bundle() {

  gutil.log('Compiling JS...');

  return bundler.bundle()
    .on('error', function(err) {
      gutil.log(err.message);
      this.emit("end");
    })
    //.pipe(exorcist('public/js/dist/bundle.js.map'))
    .pipe(source('bundle.js'))
    .pipe(gulp.dest('./server/public/js/dist/app'));
}

/**
 * Gulp task alias
 */
gulp.task('bundle', function() {
  return bundle();
});

/**
 * First bundle
 */
gulp.task('default', ['bundle']);
