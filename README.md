# wording

Create bite-sized custom Wordle-style puzzles to share with your
friends.

# Requirements

- [x] The top-level/index page displays a count of:
  - [x] Games created
  - [x] Guesses cast
  - [x] Correct guesses
- [x] A visitor can create a game from the top-level/index page.
- [x] A game creator can create a game and share a link to the puzzle they
      created.
- [x] A game creator is given a secret management URL at time of create so
      they can manage the game.
- [ ] A game creator can specify an expiration date for the game (up to a
      defined limit.)
- [x] A game creator can specify the number of attempts players will have for
      the game (up to a defined maximum.)
- [x] A game creator cannot modify guess limits for a game that has already
      been created.
- [x] A game creator can permanently delete the game with their management URL.
- [ ] A game creator can choose whether or not the answer is displayed to a
      player who has exhausted all of their guess attempts.
- [ ] A game creator can specify if game stats should be displayed.
- [ ] The game collects the following statistics for each game:
  - [ ] Total attempts
  - [ ] Number of correct guesses
  - [ ] The time it took a player to guess correctly
- [x] The UI presents visual feedback for guesses:
  - [x] Gray means a letter is not included in the word.
  - [x] Yellow means a letter is included in the word but is in the wrong
  spot.
  - [x] Green means the letter is included in the word and is in the correct
position.
- [x] A player is not able to submit a guess if:
  - [x] they have already guessed the correct answer;
  - [x] or they have no votes remaining.
- [x] A player's participation in a game is tracked and shown when they visit
      the puzzle.
- [x] A guess is not submitted if:
  - [x] the guess is not completely alphabetical (no whitespace, no numbers,
        no symbols);
  - [x] the player has already guessed that word;
  - [x] the guess is not the same length as the answer;
  - [x] the player has no guesses remaining;
- [ ] The UI offers feedback for input validation:
  - [ ] guess length must be exactly the length of the answer;
  - [ ] remaining attempts are disabled when the game is over;
  - [ ] the game is disabled when past the expiry time;
- [ ] Expired games are marked as deleted after an interval passes
