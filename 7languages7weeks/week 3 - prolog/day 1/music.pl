artist(belab,drums,punk).
artist(farin,guitar,punk).
artist(hatfield,guitar,metal).

artist_by_instrument(Who,What) :- artist(Who,What,_).
artist_by_genre(Who,What) :- artist(Who,_,What).