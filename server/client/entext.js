
var mod = angular.module('entext', ['ngSanitize']);

mod.controller('MainCtrl', function($scope, $sce, $sanitize, $http, $location) {
  $scope.response = {
  "results": [],
  "took": "0ms",
  "error": "",
  "info": {
    "e": "demo",
    "patterns": "?",
  }
}

  $scope.create_message = "Your dictionary will be assigned an unique URL. It will be forgotten if incative for one hour.";
  $scope.pattern_staging = '';
  $scope.create_patterns = [];

  var t;
  $scope.$watch('pattern_staging', function() {
    clearTimeout(t)
    t = setTimeout(function() {
      if ($scope.pattern_staging.length > 0) {
        $scope.create_patterns = $scope.pattern_staging.split('\n');
      } else {
        $scope.create_patterns = [];
      }
    }, 100)

  })

$scope.create = function() {
  if ($scope.create_patterns.length > 10000) {
      $scope.create_message = "Too many patterns";
      return
  }
  if ($scope.create_patterns.length == 0) {
      $scope.create_message = "No patterns";
      return
  }

  if ($scope.pattern_staging.length > 1025*512) {
      $scope.create_message = "Too much text";
      return
  }

  $scope.create_message = "Creating...";
  $http.post('create?k=2', $scope.create_patterns).then(function(res) {
    $location.path(res.data);
    $scope.cancel_create();
    $scope.search();
  })
}

$scope.start_create = function() {
  document.getElementById('create').style.display='block';
}

$scope.cancel_create = function() {
    $scope.create_message = "Your dictionary will be assigned an unique URL. It will be forgotten if incative for one hour.";
    document.getElementById('create').style.display='none';
    $scope.pattern_staging = '';
    $scope.create_patterns = [];
}

$scope.clear = function() {

  document.getElementById('left-context').innerHTML = '';
  document.getElementById('right-context').innerHTML = '';

  var e = $location.path().substring(1);

  if (e.length === 0) {
    e = 'demo'
  }

  if ($scope.response.info == undefined) {
    $scope.response.info = {};
  }

  $scope.response.info.e = e;
  //$scope.response.info.patterns = "?";
}

$scope.search = function(preserve) {
  $scope.clear();

  if (!preserve) {
    $scope.text = $sanitize(document.getElementById('text').innerHTML);
    var regex = new RegExp("<span([^>]*)>([^<]+)</span>", "g")
    $scope.text = $scope.text.replace(regex, "$2");
  }

  $http.post('search?e=' +  $scope.response.info.e + '&k=' + $scope.k, $scope.text).then(function(res) {

    var refs = [];

    var staging_text = $scope.text;

    $scope.response = res.data;

    if ($scope.response.results == undefined) {
      return;
    }

    $scope.response.results.sort(function(a, b) {
      if (a.Offset != b.Offset) {
        return a.Offset - b.Offset
      }

      return a.Distance - b.Distance;
    })

    var po = -100

    $scope.response.results.map(function(result, i) {
      if (Math.abs(result.Offset - po) < result.Actual.length) {
        return;
      }
      po = result.Offset;
      var regex = new RegExp(result.Actual, "g")
      var id = "hit-" + i;
      staging_text = staging_text.replace(regex, "<span class='yellow-" + result.Distance + "' id='" + id + "'>" + result.Actual +"</span>")
      refs.push({id: id, reference: result.Reference, distance: result.Distance})
    })

    $scope.staging_text = $sce.trustAsHtml(staging_text);

    setTimeout(function() {
      var cache = {'left': {}, 'right': {}}

      refs.map(function(ref) {
        elem = document.getElementById(ref.id)
        if (!elem) return
        rect = elem.getBoundingClientRect();
        pos =  window.scrollY + rect.top;
        left = rect.left < document.documentElement.clientWidth / 2

        side = left ? 'left' : 'right'

        if (cache[side][pos] == undefined) {
          e = angular.element('<div class="context-row"></div>');
          e[0].innerHTML = ref.reference + " (" + ref.distance + ")";

          e[0].setAttribute('style', 'position:absolute;top:' + pos + 'px;width:160px;');
          e[0].setAttribute('title', ref.reference);
          container = document.getElementById(side + '-context');
          container.appendChild(e[0])
          cache[side][pos] = e
        } else {
          e = cache[side][pos]
          e[0].innerHTML +=  ", " + ref.reference + " (" + ref.distance + ")";
          e[0].setAttribute('title', e[0].getAttribute('title') + ", " + ref.reference);
        }
      })
    });
  })
}

  $scope.k = 2;

  $scope.clear();



  if ($scope.response.info.e == 'demo') {
    $scope.text = "Hi there!\n\nThis is a simple demo application of an algorithm for multi-pattern approximate search. It lets you search full text for known patterns like names and organiztions, with errors (insert, edit, delete and transpose).\n\nFor example, Vasili Pushkin may be a different transliteration of Vasily Pushkin the poet, and Ryen Renolds may just be lazy writing. Think of it as a combination of Aho-Corasick and Levenshtein automata, although that's not quite how it works. "
     + "A research paper on the algorithm is currently in peer-review, and we will publish more details later.\n\nThis demo searches for more than 500 thousand names collected from Wikipedia, but you can create your own dictionary by clicking the button in the top right corner. The application searches for whatever is in this white box, so <b>click here to edit</b>, and tap 'search' below! You can also drag-n-drop a .docx, .xlsx, .pptx, .pdf, and .msg file.\n\n"
     + "Nothing you send will be stored on disk, and inactive dictionaries will be discarded. However, we do gather some metrics to measure performance. We make no guarantee of availability, confidentiality, or anything else. We'd appreciate your feedback and bug reports!\n\nGet in touch: <a href=\"mailto:demo@entext.io\">demo@entext.io";
   } else {
     $scope.text = '';
   }


   $scope.staging_text = ""

   $scope.$watch("k", function() {
     $scope.search();
   });

  var w;
  $scope.$on("$locationChangeSuccess", function() {
    clearTimeout(w)
    w = setTimeout(function() {
    $scope.clear();

    if ($scope.response.info.e == 'demo') {
      $scope.text = "Hi there!\n\nThis is a simple demo application of an algorithm for multi-pattern approximate search. It lets you search full text for known patterns like names and organiztions, with errors (insert, edit, delete and transpose).\n\nFor example, Vasili Pushkin may be a different transliteration of Vasily Pushkin the poet, and Ryen Renolds may just be lazy writing. Think of it as a combination of Aho-Corasick and Levenshtein automata, although that's not quite how it works. "
       + "A research paper on the algorithm is currently in peer-review, and we will publish more details later.\n\nThis demo searches for " + $scope.response.info.patterns + " names collected from Wikipedia, but you can create your own dictionary by clicking the button in the top right corner. The application searches for whatever is in this white box, so <b>click here to edit</b>, and tap 'search' below! You can also drag-n-drop a .docx, .xlsx, .pptx, .pdf, and .msg file.\n\n"
       + "Nothing you send will be stored on disk, and inactive dictionaries will be discarded. However, we do gather some metrics to measure performance. We make no guarantee of availability, confidentiality, or anything else. We'd appreciate your feedback and bug reports!\n\nGet in touch: <a href=\"mailto:demo@entext.io\">demo@entext.io";
     } else {
       $scope.text = '';
     }

   $scope.staging_text = $sce.trustAsHtml($scope.text);

    $scope.search(true);
    }, 10)
  })
});
