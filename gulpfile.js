var gulp = require('gulp');
var gutil = require('gulp-util');
var source = require('vinyl-source-stream');
var babelify = require('babelify');
var watchify = require('watchify');
var exorcist = require('exorcist');
var browserify = require('browserify');
var child      = require('child_process');
var sync       = require('gulp-sync')(gulp).sync;

/* ----------------------------------------------------------------------------
 * Application server
 * ------------------------------------------------------------------------- */

/*
 * Build application server.
 */
gulp.task('server:build', function() {
  var build = child.spawnSync('go', ['build'], {cwd:'server'});
  if (build.stderr.length) {
    var lines = build.stderr.toString()
      .split('\n').filter(function(line) {
        return line.length
      });
    for (var l in lines)
      gutil.log(gutil.colors.red(
        'Error (go install): ' + lines[l]
      ));
    notifier.notify({
      title: 'Error (go install)',
      message: lines
    });
  }
  return build;
});
var server = undefined;
/*
 * Restart application server.
 */
gulp.task('server:spawn', function() {
  if (server != undefined)
    server.kill();

  /* Spawn application server */
  server = child.spawn('server.exe', [], {cwd:'server'});

  /* Pretty print server log output */
  server.stdout.on('data', function(data) {
    var lines = data.toString().split('\n')
    for (var l in lines)
      if (lines[l].length)
        gutil.log(lines[l]);
  });

  /* Print errors to stdout */
  server.stderr.on('data', function(data) {
    process.stdout.write(data.toString());
  });
});

/*
 * Watch source for changes and restart application server.
 */
 gulp.task('server:watch', function() {
   /* Rebuild and restart application server */
   gulp.watch([
     '*/**/*.go',
   ], sync([
     'server:build',
     'server:spawn'
   ], 'server'));
 });



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
    .pipe(gulp.dest('.server/public/js/dist/app'));
}

/**
 * Gulp task alias
 */
gulp.task('bundle', function() {
  return bundle();
});

gulp.task('build', [
  'server:build'
]);
gulp.task('watch', [
  'server:build'
], function() {
  return gulp.start([
    'server:watch',
    'server:spawn'
  ]);
});

gulp.task('default', ['bundle', 'build']);
