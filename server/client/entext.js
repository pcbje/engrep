
var mod = angular.module('entext', ['ngSanitize']);

mod.controller('MainCtrl', function($scope, $sce) {
  $scope.response = {
  "results": [
    {
      "Actual": "Vasili Pushkin",
      "Reference": "Vasily Pushkin",
      "Distance": 1,
      "Info": ""
    },
    {
      "Actual": "Vasily Pushkin",
      "Reference": "Vasily Pushkin",
      "Distance": 0,
      "Info": ""
    },
    {
      "Actual": "Ryen Renolds",
      "Reference": "Ryan Reynolds",
      "Distance": 2,
      "Info": ""
    }
  ],
  "took": "143.192194ms",
  "info": {
    "e": "demo",
    "patterns": 569567
  }
}

$scope.k = 2;

  $scope.text = "Hi there!\n\nThis is a demo application of an algorithm for multi-pattern approximate search. It lets you search full text for known patterns like names and organiztions, with errors (insert, edit, delete and transpose).\n\nFor example, Vasili Pushkin may be a different transliteration of Vasily Pushkin the poet, and Ryen Renolds may just be lazy writing. Think of it as a combination of Aho-Corasick and Levenshtein automata, although that's not quite how it works."
   + "A research paper on the algorithm is currently in peer-review, and we will publish more details later.\n\nThis demo searches for more than 500 thousand names collected from Wikipedia, but you can create your own dictionary by clicking the button in the top right corner. The application searches for whatever is in this white box, so <b>click here to edit</b>, and tap 'search' below! You can also drag-n-drop a .docx, .xlsx, .pptx, .pdf, and .msg file.\n\n"
   + "Nothing you send will be stored on disk, and inactive dictionaries will be discarded. However, we do gather some metrics to measure performance. We make no guarantee of availability, confidentiality, or anything else, and we'd appreciate your feedback and bug reports!\n\nGet in touch: <a href=\"mailto:demo@entext.io\">demo@entext.io";

  var staging_text = $scope.text;

  var refs = [];

  $scope.response.results.map(function(result, i) {
      var regex = new RegExp(result.Actual, "gi")
      var id = "hit-" + i;
      staging_text = staging_text.replace(regex, "<span id='" + id + "' class='yellow-" + result.Distance + "'>" + result.Actual +"</span>")
      refs.push({id: id, reference: result.Reference, distance: result.Distance})
  })

  $scope.staging_text = $sce.trustAsHtml(staging_text);


  setTimeout(function() {
    var cache = {'left': {}, 'right': {}}
    refs.map(function(ref) {
      elem = document.getElementById(ref.id)
      rect = elem.getBoundingClientRect();
      pos = rect.top
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
        e.innerHTML +=  ", " + ref.reference + " (" + ref.distance + ")";
        e[0].setAttribute('title', e[0].getAttribute('title') + ", " + ref.reference);
      }



    })
  });

});
