# wording

Create bite-sized custom Wordle-style puzzles to share with your
friends.

# Requirements

- [ ] The top-level/index page can display a count of:
  - [ ] Games created
  - [ ] Guesses cast
  - [ ] Correct guesses
- [ ] A visitor can create a game from the top-level/index page.
- [ ] A game creator can create a game and share a link to the puzzle they
      created.
- [ ] A game creator is given a secret management URL at time of create so
      they can manage the game.
- [ ] A game creator can specify an expiration date for the game (up to a
      defined limit.)
- [ ] A game creator can specify the number of attempts players will have for
      the game (up to a defined maximum.)
- [ ] A game creator can expire the game early with their management URL.
- [ ] A game creator can un-expire a game that they accidentally marked as
      expired with their management URL.
- [ ] A game creator cannot modify guess limits for a game that has already
      been created.
- [ ] A game creator can permanently delete the game with their management URL.
- [ ] A game creator can choose whether or not the answer is displayed to a
      player who has exhausted all of their guess attempts.
- [ ] A game creator can specify if game stats should be displayed.
- [ ] The game collects the following statistics for each game:
  - [ ] Total attempts
  - [ ] Number of correct guesses
  - [ ] The time it took a player to guess correctly
- [ ] A player can submit a guess so long as the amount of previous guesses
      does not exceed the maximum allowed for the puzzle.
- [ ] The UI presents visual feedback for guesses:
  - [ ] Gray means a letter is not included in the word.
  - [ ] Yellow means a letter is included in the word but is in the wrong
  spot.
  - [ ] Green means the letter is included in the word and is in the correct
position.
- [ ] A player is not able to submit a guess if:
  - they have already guessed the correct answer;
  - or they have no votes remaining.
- [ ] A player's participation in a game is tracked and shown when they visit
      the puzzle.
