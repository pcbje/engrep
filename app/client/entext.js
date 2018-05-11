'use strict';

var modules = angular.module('modules', []);
var app = angular.module('entext', ['ngSanitize', 'modules']);

modules.service('api', function() {
    return {}
})

modules.service('search', function($rootScope, $sce, $sanitize, $http, $location, api) {
  var timeout;
  var e;
  var intro = "Hi there!\n\nThis is a simple demo application of an algorithm for multi-pattern approximate search. It lets you search full text for known patterns like names and organiztions, with errors (insert, edit, delete and transpose).\n\nFor example, Vasili Pushkin may be a different transliteration of Vasily Pushkin the poet, and Ryen Renolds may just be lazy writing. Think of it as a combination of Aho-Corasick and Levenshtein automata, although that's not quite how it works. "
     + "We are working on a research paper for this algorithm, and will publish more details later.\n\nThis demo searches for 539364 names collected from Wikipedia, but you can create your own dictionary by clicking the button in the top right corner. The application searches for whatever is in this white box, so <b>click here to edit</b> and then hit 'search' below!\n\n"
     + "Nothing you send will be stored on disk, and inactive dictionaries will be discarded. However, we do gather some metrics to measure performance.\n\nQuestions or comments? Get in touch! <a href=\"mailto:demo@entext.io\">demo@entext.io</a>\n\n";

  document.getElementById('text').innerHTML = intro;

  $rootScope.$on("$locationChangeSuccess", function() {
    api.search.vars.dictionary = $location.path().substring(1);

    if (api.search.vars.dictionary.length === 0) {
      api.search.vars.dictionary = 'demo'
    }

    api.search.clear();

    api.search.vars.text = '';

    if (api.search.vars.dictionary == 'demo') {
      api.search.vars.text = intro;
    }
  });

  var regex1 = new RegExp("<span([^>]*)>", "g");
  var regex2 = new RegExp("</span>", "g");
  var prep = new RegExp(" ", "g");

  api.search = {
    vars: {
      text: intro,
      dictionary: '',
      k: 2,
      patterns: 0,
      hits: 0,
      took: 0.0,
    },
    clear: function() {
      document.getElementById('left-context').innerHTML = '';
      document.getElementById('right-context').innerHTML = '';

      api.search.vars.text = document.getElementById('text').innerHTML;
      api.search.vars.text = api.search.vars.text.replace(regex1, "");
      api.search.vars.text = api.search.vars.text.replace(regex2, "");
    },
    search: function() {
      api.search.clear();

      $http.post('search?d=' +  api.search.vars.dictionary + '&k=' + api.search.vars.k, api.search.vars.text).then(function(res) {
        var refs = [];

        api.search.vars.hits = 0;

        var response = res.data;

        api.search.vars.patterns = response.patterns;
        api.search.vars.took = response.took;

        if (response.results == undefined) {
          return;
        }

        response.results.sort(function(a, b) {
          if (a.Offset != b.Offset) {
            return a.Offset - b.Offset
          }

          return a.Distance - b.Distance;
        });

        var po = -100
        var pr = null

        response.results.map(function(result, i) {
          if (Math.abs(result.Offset - po) <= api.search.vars.k) {
            return;
          }

          api.search.vars.hits++;

          po = result.Offset;
          pr = result.Reference;

          var regex = new RegExp(result.Actual.replace(prep, "( |</span>)*"), "g")
          var id = "hit-" + i;

          api.search.vars.text = api.search.vars.text.replace(regex, "<span class='yellow-" + result.Distance + "' id='" + id + "'>" + result.Actual +"</span>")
          refs.push({id: id, reference: result.Reference, distance: result.Distance})
        })

        api.search.vars.text = $sce.trustAsHtml(api.search.vars.text);

        setTimeout(function() {
          var cache = {'left': {}, 'right': {}}

          refs.map(function(ref) {
            var elem = document.getElementById(ref.id)
            if (!elem) return
            var rect = elem.getBoundingClientRect();
            var pos =  window.scrollY + rect.top;
            var left = rect.left < document.documentElement.clientWidth / 2

            var side = left ? 'left' : 'right'
            var container = document.getElementById(side + '-context');

            if (cache[side][pos] == undefined) {
              e = angular.element('<div class="context-row"></div>');
              e[0].innerHTML = ref.reference + " (" + ref.distance + ")";

              e[0].setAttribute('style', 'position:absolute;top:' + pos + 'px;width:160px;');
              e[0].setAttribute('title', ref.reference);
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
  }
});

modules.service('dictionary', function($location, $http, api) {
  var create_message = "The dictionary will be assigned an unique URL. It will be forgotten if incative.";

  api.dictionary = {
    vars: {
      patterns: '',
      message: create_message,
    },
    start_create: function() {
      document.getElementById('create').style.display = 'block';
    },
    terminate_create: function() {
      document.getElementById('create').style.display='none';
      api.dictionary.vars.message = create_message;
      api.dictionary.vars.patterns = '';
    },
    create: function() {
      var patterns = [];
      var pattern_string = api.dictionary.vars.patterns;

      if (pattern_string.length > 1025*512) {
          api.dictionary.vars.message = "Too much text";
          return
      }

      if (pattern_string.length > 0) {
          patterns = pattern_string.split('\n');
      }

      if (patterns.length > 10000) {
          api.dictionary.vars.message = "Too many patterns";
          return
      }

      if (patterns.length == 0) {
          api.dictionary.vars.message = "No patterns";
          return
      }

      api.dictionary.vars.message = "Creating...";

      $http.post('create?k=2', patterns).then(function(res) {
        $location.path(res.data);
        api.dictionary.terminate_create();
        api.search.search();
      })
    }
  }
})

app.controller('MainCtrl', function($scope, api, dictionary, search) {
  $scope.api = api;

  $scope.$watch("api.search.vars.k", function() {
   api.search.search();
  });
});
