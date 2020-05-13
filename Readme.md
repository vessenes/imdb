# Parsing IMDB's Public Data For Fun

This repository will process locally downloaded IMDB files into a short-term in-memory database suitable for querying and navigating IMDB.

This software was built because our family movie night requires that we choose our next film based on the following two rules:

1. The next film must feature an actor featured in our most recent film
2. The actor linking the two films must not be the actor linking our two most recent films

We had planned on using this to get between any two films, but ended up using it for just mucking about and picking our next film. It does warn us of dangerous neighborhoods (films with only poor quality related films), which is helpful. 

# Disclaimers

You probably don't want this, really. 

##Features 

* CLI based
* Find Movies Related To Each-Other
* Calculate Actor Rating Value Over Replacement

## Upcoming Features / Ideas

* API?
* Six Degrees Search - Navigate between two movies using only linkage rules 
* Add Box Office Value Ratings (will need IMDB Pro data for this)

# Installation

```bash
$ go get
$ go generate    # tinylib's msgp generates custom marshalling functions
$ go build
$ ./download.sh  # This downloads gigabytes of IMDB data
$ ./imdb         # First incantation will create the graph database
$ ./imdb         # Future Incantations will load the database ready for querying
```

# Usage

The Command Line is search oriented. It uses tab completion to help you find a movie or an actor.

It responds to names of films, actors, and also IMDB unique ids (ttxxxx and nmxxxx).


## Film Lookup

```bash

> Grosse Pointe Blank

film:  Grosse Pointe Blank
id:  tt0119229
rated: 73
Broadly Rated? true
actors:
+-----------+--------------------+------+--------+
| nm0294062 | Colby French       | 6.75 |  0.302 |
| nm0603890 | Belita Moreno      | 6.77 |  0.316 |
| nm0223258 | Nicholas de Wolff  | 6.77 |  0.286 |
| nm0082711 | Laurence Bilzerian | 6.00 | -0.044 |
| nm0193639 | Ann Cusack         | 6.28 |  0.155 |
| nm0035497 | Brent Armitage     | 6.20 |  0.076 |
| nm0415070 | Carlos Jacott      | 6.85 |  0.262 |
| nm0000131 | John Cusack        | 6.14 |  0.013 |
| nm0005315 | Jeremy Piven       | 6.21 |  0.058 |
| nm0694043 | Brian Powell       | 7.30 |  0.478 |
| nm0735313 | Eva Rodriguez      | 6.83 |  0.303 |
| nm0230151 | K.K. Dodds         | 6.78 |  0.286 |
| nm0861265 | Wendy Thorlakson   | 5.85 | -0.045 |
| nm0000101 | Dan Aykroyd        | 5.98 | -0.110 |
| nm0000349 | Joan Cusack        | 6.10 | -0.036 |
| nm0458564 | Anthony Kleeman    | 5.40 | -0.036 |
| nm0552513 | Jim Martin         | 7.15 |  0.616 |
| nm0293461 | K. Todd Freeman    | 7.38 |  0.511 |
| nm0193641 | Bill Cusack        | 7.45 |  0.580 |
| nm0222586 | Sarah DeVincentis  | 6.58 |  0.229 |
| nm0000378 | Minnie Driver      | 6.24 |  0.039 |
| nm0191044 | Michael Cudlitz    | 5.40 | -0.191 |
| nm0457389 | Audrey Kissel      | 7.30 |  0.769 |
| nm0364455 | Barbara Harris     | 6.78 |  0.205 |
| nm0233723 | Traci Dority       | 7.30 |  0.769 |
| nm0752751 | Mitchell Ryan      | 6.82 |  0.241 |
+-----------+--------------------+------+--------+

```

## Actor Lookup 
```
> Joan Cusac<tab>
Actors [Joan Cusack]

Actor: Joan Cusack
ID:  nm0000349
Value: -0.0363015335880794
Average Rating: 6.103448275862069
Filmography:
+-----------+--------------------------------+--------+----------------+
|    ID     |              FILM              | RATING | BROADLY RATED? |
+-----------+--------------------------------+--------+----------------+
| tt0100212 | My Blue Heaven                 |     63 | true           |
| tt0332379 | School of Rock                 |     70 | true           |
| tt0119229 | Grosse Pointe Blank            |     73 | true           |
| tt1342403 | Toys in the Attic              |     61 | false          |
| tt0396652 | Ice Princess                   |     60 | true           |
| tt0100134 | Men Don't Leave                |     65 | true           |
| tt0113986 | Nine Months                    |     55 | true           |
| tt0106220 | Addams Family Values           |     66 | true           |
| tt0117102 | Mr. Wrong                      |     39 | true           |
| tt0119360 | In & Out                       |     64 | true           |
| tt0104412 | Hero                           |     65 | true           |
| tt0120151 | A Smile Like Yours             |     49 | true           |
| tt0163187 | Runaway Bride                  |     55 | true           |
| tt1305591 | Mars Needs Moms                |     54 | true           |
| tt0350028 | Raising Helen                  |     60 | true           |
| tt1093908 | Confessions of a Shopaholic    |     59 | true           |
| tt0105629 | Toys                           |     50 | true           |
| tt0436331 | Friends with Money             |     56 | true           |
| tt0120363 | Toy Story 2                    |     79 | true           |
| tt0092537 | The Allnighter                 |     41 | true           |
| tt0371606 | Chicken Little                 |     56 | true           |
| tt0096463 | Working Girl                   |     68 | true           |
| tt0435761 | Toy Story 3                    |     83 | true           |
| tt0150216 | Cradle Will Rock               |     68 | true           |
| tt0884224 | War, Inc.                      |     55 | true           |
| tt0101530 | The Cabinet of Dr. Ramirez     |     60 | false          |
| tt2710534 | Arrive Alive                   |     64 | false          |
| tt0844993 | Hoodwinked Too! Hood vs. Evil  |     46 | true           |
| tt2338454 | Unicorn Store                  |     55 | true           |
| tt0137363 | Arlington Road                 |     71 | true           |
| tt1659337 | The Perks of Being a           |     80 | true           |
|           | Wallflower                     |        |                |
| tt0846308 | Kit Kittredge: An American     |     65 | true           |
|           | Girl                           |        |                |
+-----------+--------------------------------+--------+----------------+
```

## Finding Linked Movies

Linkage searches for all movies that share actors with the source, and orders by IMDB rating. 

Ask for a link by entering ```l:``` and the IMDB film id. 

```bash
> Three Colors: Blue
film:  Three Colors: Blue
id:  tt0108394
rated: 79
Broadly Rated? true
actors:
+-----------+------------------------+------+--------+
| nm0904996 | Charlotte Véry         | 6.80 |  0.111 |
| nm0000300 | Juliette Binoche       | 6.49 | -0.026 |
| nm0147544 | Idit Cebula            | 7.90 |  0.554 |
| nm0952498 | Zbigniew Zamachowski   | 7.21 |  0.176 |
| nm0286248 | Pierre Forget          | 7.67 |  0.309 |
| nm0674340 | Florence Pernel        | 7.20 |  0.236 |
| nm0898653 | Hélène Vincent         | 7.11 |  0.137 |
| nm0703364 | Hugues Quester         | 6.92 |  0.082 |
| nm0728938 | Emmanuelle Riva        | 7.33 |  0.270 |
| nm0753666 | Benoît Régent          | 8.00 |  0.582 |
| nm0242211 | Claude Duneton         | 7.70 |  0.317 |
| nm0635048 | Stanislas Nordey       | 7.90 |  0.554 |
| nm0671615 | Yves Penay             | 7.20 |  0.262 |
| nm0755652 | Isabelle Sadoyan       | 7.16 |  0.139 |
| nm0552134 | Daniel Martin          | 7.33 |  0.414 |
| nm0000365 | Julie Delpy            | 6.47 | -0.038 |
| nm0542011 | Philippe Manesse       | 7.90 |  0.554 |
| nm0857616 | Catherine Therouenne   | 7.90 |  0.554 |
| nm0647366 | Alain Ollivier         | 7.90 |  0.438 |
| nm0874556 | Yann Trégouët          | 7.90 |  0.554 |
| nm0156756 | Arno Chevrier          | 5.80 | -0.024 |
| nm0605429 | Philippe Morier-Genoud | 7.75 |  0.332 |
| nm0901812 | Philippe Volter        | 7.73 |  0.365 |
| nm0228432 | Jacques Disses         | 7.60 |  0.314 |
+-----------+------------------------+------+--------+

> l:tt0108394

Finding all linked films, connected by major actors, from tt0108394
2020/05/13 11:16:41 tt0108394
+------------+--------------------------------+------------------------+--------+----------------+
|  FILM ID   |              FILM              |         LINKER         | RATING | BROADLY RATED? |
+------------+--------------------------------+------------------------+--------+----------------+
| tt0253999  | L'île aux oiseaux              | Benoît Régent          |      0 | false          |
| tt0188275  | Visage de chien                | Hugues Quester         |      0 | false          |
| tt0342071  | L'avenir de Jéremy             | Philippe Manesse       |      0 | false          |
| tt7780996  | Radio nuit Paris               | Philippe Morier-Genoud |      0 | false          |
| tt0149237  | Tant pis si je meurs           | Philippe Manesse       |      0 | false          |
| tt0110921  | Péché véniel... péché          | Isabelle Sadoyan       |      0 | false          |
|            | mortel...                      |                        |        |                |
| tt2781994  | Love Love Love                 | Arno Chevrier          |      0 | false          |
| tt0885510  | L'or et le plomb               | Emmanuelle Riva        |      0 | false          |
| tt0188966  | Una notte, un sogno            | Hugues Quester         |      0 | false          |
| tt0104970  | Naprawde krotki film o         | Zbigniew Zamachowski   |      0 | false          |
|            | milosci, zabijaniu i jeszcze   |                        |        |                |
|            | jednym przykazaniu             |                        |        |                |
| tt0236578  | Paroles d'hommes               | Florence Pernel        |      0 | false          |
| tt0108678  | L'écrivain public              | Florence Pernel        |      0 | false          |
| tt0103188  | La valse des pigeons           | Arno Chevrier          |      0 | false          |
| tt11368866 | Tarapaty 2                     | Zbigniew Zamachowski   |      0 | false          |
| tt0097853  | Les matins chagrins            | Hugues Quester         |      0 | false          |
| tt0341214  | Accord parfait                 | Benoît Régent          |      0 | false          |
| tt0309586  | La femme intégrale             | Benoît Régent          |      0 | false          |
| tt8767704  | La malaimée                    | Arno Chevrier          |      0 | false          |
| tt10310144 | Camino Real                    | Juliette Binoche       |      0 | false          |
| tt6100508  | Alma                           | Emmanuelle Riva        |      0 | false          |
| tt5850650  | Together Now                   | Juliette Binoche       |      0 | false          |
| tt0259949  | The Whip                       | Pierre Forget          |      0 | false          |
| tt7438706  | Adorables                      | Hélène Vincent         |      0 | false          |
| tt1291550  | La mort de l'utopie            | Emmanuelle Riva        |      0 | false          |
| tt0119036  | Le déménagement                | Yann Trégouët          |      0 | false          |
| tt0123631  | Capitaine au long cours        | Emmanuelle Riva        |      0 | false          |
| tt0252246  | Attendre le navire             | Benoît Régent          |      0 | false          |
| tt10088984 | Between Two Worlds             | Juliette Binoche       |      0 | false          |
| tt1623754  | La mémoire de l'eau            | Florence Pernel        |      0 | false          |
| tt10551904 | L'origine du monde             | Hélène Vincent         |      0 | false          |
| tt0284536  | Le syndrome de Peter Pan       | Charlotte Véry         |      0 | false          |
| tt0086888  | My Friend Washington           | Benoît Régent          |     25 | false          |
| tt1708498  | Projekt dziecko, czyli ojciec  | Zbigniew Zamachowski   |     27 | false          |
|            | potrzebny od zaraz             |                        |        |                |
| tt0093673  | On a volé Charlie Spencer!     | Stanislas Nordey       |     30 | false          |
| tt0418287  | Le veilleur                    | Arno Chevrier          |     32 | false          |
| tt9643428  | Armani Privé - A view beyond   | Juliette Binoche       |     37 | false          |
| tt0402294  | Iznogoud- Caliph Instead of    | Arno Chevrier          |     37 | true           |
|            | the Caliph                     |                        |        |                |
| tt0150331  | Les diplômés du dernier rang   | Philippe Manesse       |     39 | false          |
| tt5338174  | Les naufragés                  | Philippe Morier-Genoud |     40 | false          |
| tt0120871  | The Treat                      | Julie Delpy            |     40 | false          |
| tt0099477  | Dédé                           | Hélène Vincent         |     40 | false          |
| tt0105693  | Une journée chez ma mère       | Hélène Vincent         |     41 | false          |
| tt0098514  | Trois années                   | Philippe Volter        |     43 | false          |
| tt0089613  | My Brother-in-law Killed My    | Juliette Binoche       |     44 | false          |
|            | Sister                         |                        |        |                |
| tt0100125  | Mauvaise fille                 | Florence Pernel        |     44 | false          |
| tt0317535  | Fureur                         | Yann Trégouët          |     45 | false          |
| tt0258685  | Issue de secours               | Philippe Volter        |     45 | false          |
| tt0243991  | Intimate Affairs               | Julie Delpy            |     46 | true           |
| tt0130201  | Pulapka                        | Zbigniew Zamachowski   |     46 | false          |
| tt5312370  | 7 rzeczy, których nie wiecie o | Zbigniew Zamachowski   |     46 | false          |
|            | facetach                       |                        |        |                |
| tt2391746  | Rue Mandar                     | Idit Cebula            |     48 | false          |
| tt0428765  | The Legend of Lucy Keyes       | Julie Delpy            |     48 | true           |
| tt0488928  | Guilty Hearts                  | Julie Delpy            |     48 | false          |
| tt4357764  | No Panic, With a Hint of       | Zbigniew Zamachowski   |     50 | false          |
|            | Hysteria                       |                        |        |                |
| tt0119674  | Les mille merveilles de        | Julie Delpy            |     50 | false          |
|            | l'univers                      |                        |        |                |
| tt1691826  | Ariane                         | Emmanuelle Riva        |     50 | false          |
| tt1535612  | The Son of No One              | Juliette Binoche       |     50 | true           |
| tt1480656  | Cosmopolis                     | Juliette Binoche       |     50 | true           |
| tt5975354  | Baby Bump(S)                   | Juliette Binoche       |     50 | true           |
| tt0109929  | Grande petite                  | Hugues Quester         |     50 | false          |
| tt0118604  | An American Werewolf in Paris  | Julie Delpy            |     50 | true           |
| tt0399757  | Thierry Mugler                 | Juliette Binoche       |     50 | false          |
| tt5338644  | Mrs. Hyde                      | Charlotte Véry         |     51 | true           |
| tt0155964  | Gates of Fire                  | Emmanuelle Riva        |     51 | false          |
| tt3686942  | Killing Love                   | Zbigniew Zamachowski   |     51 | false          |
| tt0108357  | Tom est tout seul              | Hélène Vincent         |     53 | false          |
| tt0306011  | Resistance                     | Philippe Volter        |     53 | false          |
| tt7566518  | Vision                         | Juliette Binoche       |     53 | false          |
| tt0490931  | Les filles de Malemort         | Pierre Forget          |     53 | false          |
| tt0446442  | A Few Days in September        | Juliette Binoche       |     53 | true           |
| tt9568486  | Mine de rien                   | Hélène Vincent         |     54 | false          |
| tt6692840  | Du soleil dans mes yeux        | Hélène Vincent         |     54 | false          |
| tt0292964  | Beginner's Luck                | Julie Delpy            |     55 | false          |
| tt0103611  | Abracadabra                    | Philippe Volter        |     55 | false          |
| tt7217028  | Joint Custody                  | Hélène Vincent         |     55 | false          |
| tt1549589  | Elles                          | Juliette Binoche       |     55 | true           |
| tt0179206  | War in the Highlands           | Yann Trégouët          |     55 | false          |
| tt0977647  | Deux vies... plus une          | Idit Cebula            |     55 | false          |
| tt6081632  | Marie-Francine                 | Hélène Vincent         |     55 | false          |
| tt0387059  | Bee Season                     | Juliette Binoche       |     55 | true           |
| tt2088962  | An Open Heart                  | Juliette Binoche       |     55 | false          |
| tt0093349  | King Lear                      | Julie Delpy            |     56 | true           |
| tt1860260  | Blood from a Stone             | Yann Trégouët          |     56 | false          |
| tt0871512  | Disengagement                  | Juliette Binoche       |     56 | false          |
| tt4085944  | Lolo                           | Julie Delpy            |     56 | true           |
| tt0385363  | Let's Make a Grandson          | Zbigniew Zamachowski   |     56 | false          |
| tt0076473  | Golden Night                   | Catherine Therouenne   |     56 | false          |
| tt0210816  | Looking for Jimmy              | Julie Delpy            |     56 | false          |
| tt4827558  | High Life                      | Juliette Binoche       |     58 | true           |
| tt0096965  | Les bois noirs                 | Philippe Volter        |     58 | false          |
| tt0425236  | Mary                           | Juliette Binoche       |     58 | true           |
| tt0096879  | Bal na dworcu w Koluszkach     | Zbigniew Zamachowski   |     58 | false          |
| tt0096041  | Savannah                       | Daniel Martin          |     58 | false          |
| tt0118687  | Le bassin de J.W.              | Hugues Quester         |     58 | false          |
| tt0105522  | Tak tak                        | Zbigniew Zamachowski   |     58 | false          |
| tt5114982  | A Bun in the Oven              | Hélène Vincent         |     59 | false          |
| tt4144190  | Wiener-Dog                     | Julie Delpy            |     59 | true           |
| tt0109600  | Dernier stade                  | Philippe Volter        |     59 | false          |
| tt1654829  | Thérèse                        | Isabelle Sadoyan       |     60 | true           |
| tt3150862  | Gabriel                        | Zbigniew Zamachowski   |     60 | false          |
| tt0108182  | Souvenir                       | Hugues Quester         |     60 | false          |
| tt4726636  | Slack Bay                      | Juliette Binoche       |     60 | true           |
| tt0349260  | In My Country                  | Juliette Binoche       |     60 | true           |
| tt6423776  | Let the Sunshine In            | Juliette Binoche       |     60 | true           |
| tt0293116  | Jet Lag                        | Juliette Binoche       |     60 | true           |
| tt0476755  | Czas surferów                  | Zbigniew Zamachowski   |     60 | false          |
| tt0209464  | Villa des roses                | Julie Delpy            |     60 | false          |
| tt0106531  | La cavale des fous             | Florence Pernel        |     60 | false          |
| tt6313378  | Memoir of War                  | Stanislas Nordey       |     60 | false          |
| tt0073196  | Je t'aime moi non plus         | Hugues Quester         |     60 | true           |
| tt0092426  | Pierscien i róza               | Zbigniew Zamachowski   |     60 | false          |
| tt4250606  | Good Luck Algeria              | Hélène Vincent         |     60 | false          |
| tt1602472  | 2 Days in New York             | Julie Delpy            |     60 | true           |
| tt0295430  | Les matous sont romantiques    | Philippe Manesse       |     60 | false          |
| tt0120261  | Szczesliwego Nowego Jorku      | Zbigniew Zamachowski   |     60 | false          |
| tt0118018  | A Couch in New York            | Juliette Binoche       |     60 | true           |
| tt0242154  | Yoyes                          | Florence Pernel        |     60 | false          |
| tt2936884  | Lost in Paris                  | Emmanuelle Riva        |     61 | true           |
| tt0117229  | Odwiedz mnie we snie           | Zbigniew Zamachowski   |     61 | false          |
| tt0118003  | Tykho Moon                     | Julie Delpy            |     61 | false          |
| tt2735996  | Endless Night                  | Juliette Binoche       |     61 | true           |
| tt6290584  | My Zoe                         | Julie Delpy            |     61 | false          |
| tt0094690  | L'autre nuit                   | Julie Delpy            |     61 | false          |
| tt1711484  | The Conquest                   | Florence Pernel        |     61 | true           |
| tt0126004  | The Iron Rose                  | Hugues Quester         |     61 | true           |
| tt6899268  | Territory of Love              | Daniel Martin          |     61 | false          |
| tt0442207  | Locked Out                     | Hélène Vincent         |     61 | true           |
| tt1629251  | Young Girls in Black           | Isabelle Sadoyan       |     61 | false          |
| tt0228750  | Proof of Life                  | Zbigniew Zamachowski   |     61 | true           |
| tt0496634  | The Countess                   | Julie Delpy            |     61 | true           |
| tt0996956  | Lady Jane                      | Yann Trégouët          |     61 | false          |
| tt0085432  | Le destin de Juliette          | Pierre Forget          |     61 | false          |
| tt0094002  | Keep Your Right Up             | Isabelle Sadoyan       |     61 | false          |
| tt10441958 | How to Be a Good Wife          | Juliette Binoche       |     63 | false          |
| tt0105342  | L'ombre du doute               | Emmanuelle Riva        |     63 | false          |
| tt0377651  | Cialo                          | Zbigniew Zamachowski   |     63 | false          |
| tt0084432  | The Eyes, the Mouth            | Emmanuelle Riva        |     63 | false          |
| tt0065385  | The Wedding Ring               | Isabelle Sadoyan       |     63 | false          |
| tt0094218  | A Flame in My Heart            | Benoît Régent          |     63 | false          |
| tt0144164  | Demons of War                  | Zbigniew Zamachowski   |     63 | false          |
| tt0176422  | Alice and Martin               | Juliette Binoche       |     63 | true           |
| tt0258897  | Prymas. Trzy lata z tysiaca    | Zbigniew Zamachowski   |     63 | false          |
| tt0113997  | Black for Remembrance          | Benoît Régent          |     63 | false          |
| tt0343513  | Vert paradis                   | Emmanuelle Riva        |     63 | false          |
| tt1219827  | Ghost in the Shell             | Juliette Binoche       |     63 | true           |
| tt1817191  | Another Woman's Life           | Juliette Binoche       |     63 | true           |
| tt8323120  | The Truth                      | Juliette Binoche       |     64 | true           |
| tt1988695  | The Secret of the Ant-Child    | Yann Trégouët          |     64 | false          |
| tt0177746  | The Children of the Century    | Juliette Binoche       |     64 | true           |
| tt0105783  | Warsaw: Year 5703              | Julie Delpy            |     64 | false          |
| tt1638350  | Skylab                         | Julie Delpy            |     64 | true           |
| tt0120455  | Violetta, the Motorcycle Queen | Florence Pernel        |     64 | false          |
| tt4271910  | Everything Now                 | Idit Cebula            |     64 | false          |
| tt0097292  | Es ist nicht leicht ein Gott   | Hugues Quester         |     64 | false          |
|            | zu sein                        |                        |        |                |
| tt0093709  | Beatrice                       | Julie Delpy            |     65 | false          |
| tt0113909  | Don't Forget You're Going to   | Stanislas Nordey       |     65 | false          |
|            | Die                            |                        |        |                |
| tt1509638  | Wandering Streams              | Hélène Vincent         |     65 | false          |
| tt0240009  | Le soleil au-dessus des nuages | Hélène Vincent         |     65 | false          |
| tt3715122  | L'attesa                       | Juliette Binoche       |     65 | true           |
| tt0089693  | No Man's Land                  | Hugues Quester         |     65 | false          |
| tt0411118  | Anthony Zimmer                 | Yves Penay             |     65 | true           |
| tt0154551  | Le gros coup                   | Emmanuelle Riva        |     65 | false          |
| tt0380832  | Zurek                          | Zbigniew Zamachowski   |     65 | false          |
| tt1160015  | Born in 68                     | Yann Trégouët          |     65 | false          |
| tt0092131  | Once Around the Park           | Juliette Binoche       |     65 | false          |
| tt0087144  | Dangerous Moves                | Benoît Régent          |     65 | false          |
| tt8307814  | Pól wieku poezji pózniej       | Zbigniew Zamachowski   |     65 | false          |
| tt0102930  | A Mere Mortal                  | Philippe Volter        |     65 | false          |
| tt0107788  | La veillée                     | Idit Cebula            |     65 | false          |
| tt2018086  | Camille Claudel 1915           | Juliette Binoche       |     65 | true           |
| tt0268509  | Cabin Fever                    | Zbigniew Zamachowski   |     65 | false          |
| tt0072873  | The Devil in the Heart         | Emmanuelle Riva        |     65 | false          |
| tt0103646  | Aline                          | Philippe Volter        |     65 | false          |
| tt0095569  | La maison de Jeanne            | Benoît Régent          |     65 | false          |
| tt0826711  | Flight of the Red Balloon      | Juliette Binoche       |     65 | true           |
| tt4302562  | The Scent of Mandarin          | Hélène Vincent         |     65 | false          |
| tt0217116  | La vache et le président       | Florence Pernel        |     65 | false          |
| tt0354808  | The Birch-Tree Meadow          | Zbigniew Zamachowski   |     65 | false          |
| tt2380331  | Words and Pictures             | Juliette Binoche       |     65 | true           |
| tt7250056  | Non-Fiction                    | Juliette Binoche       |     65 | true           |
| tt2113820  | Walesa: Man of Hope            | Zbigniew Zamachowski   |     65 | true           |
| tt0089366  | Hail Mary                      | Juliette Binoche       |     65 | true           |
| tt0089902  | Rendez-vous                    | Juliette Binoche       |     65 | true           |
| tt0110265  | Killing Zoe                    | Julie Delpy            |     65 | true           |
| tt0058650  | Thomas the Impostor            | Emmanuelle Riva        |     65 | false          |
| tt0086550  | Vive la sociale!               | Catherine Therouenne   |     65 | false          |
| tt0487490  | Death's Glamour                | Charlotte Véry         |     65 | false          |
| tt0102050  | Voyager                        | Julie Delpy            |     66 | true           |
| tt0216684  | Man of Desire                  | Emmanuelle Riva        |     66 | false          |
| tt0104514  | L'instinct de l'ange           | Hélène Vincent         |     66 | false          |
| tt4383288  | Polina, danser sa vie          | Juliette Binoche       |     66 | true           |
| tt0480242  | Dan in Real Life               | Juliette Binoche       |     66 | true           |
| tt2452254  | Clouds of Sils Maria           | Juliette Binoche       |     66 | true           |
| tt0057390  | The Hours of Love              | Emmanuelle Riva        |     66 | false          |
| tt0379063  | Squint Your Eyes               | Zbigniew Zamachowski   |     66 | false          |
| tt0841044  | 2 Days in Paris                | Julie Delpy            |     68 | true           |
| tt1242517  | Father, Son & Holy Cow         | Zbigniew Zamachowski   |     68 | false          |
| tt3027648  | Stacja Warszawa                | Zbigniew Zamachowski   |     68 | false          |
| tt0230921  | Family Pack                    | Hélène Vincent         |     68 | false          |
| tt7552686  | Who You Think I Am             | Juliette Binoche       |     68 | true           |
| tt0104181  | Wuthering Heights              | Juliette Binoche       |     68 | true           |
| tt0260198  | C'est la vie                   | Emmanuelle Riva        |     68 | false          |
| tt0055854  | Climats                        | Emmanuelle Riva        |     68 | false          |
| tt0189028  | House Arrest                   | Hélène Vincent         |     68 | false          |
| tt0104237  | Damage                         | Juliette Binoche       |     68 | true           |
| tt1242521  | Rebellion                      | Daniel Martin          |     69 | true           |
| tt0256787  | The Boscop Diagram             | Philippe Manesse       |     69 | false          |
| tt0088354  | Family Life                    | Juliette Binoche       |     69 | false          |
| tt0106363  | Le bateau de mariage           | Florence Pernel        |     69 | false          |
| tt0102137  | J'entends plus la guitare      | Benoît Régent          |     69 | false          |
| tt0070235  | I Will Walk Like a Crazy Horse | Emmanuelle Riva        |     69 | false          |
| tt0085839  | Liberté, la nuit               | Emmanuelle Riva        |     69 | false          |
| tt0102136  | I Don't Kiss                   | Hélène Vincent         |     69 | true           |
| tt0239234  | The Lady and the Duke          | Charlotte Véry         |     69 | true           |
| tt1082009  | Queen to Play                  | Daniel Martin          |     69 | true           |
| tt1899270  | A Few Hours of Spring          | Hélène Vincent         |     69 | false          |
| tt0096386  | Life Is a Long Quiet River     | Hélène Vincent         |     69 | true           |
| tt2006295  | The 33                         | Juliette Binoche       |     69 | true           |
| tt0062204  | Risky Business                 | Emmanuelle Riva        |     70 | false          |
| tt0288491  | Hi, Tereska                    | Zbigniew Zamachowski   |     70 | true           |
| tt0114204  | Pulkownik Kwiatkowski          | Zbigniew Zamachowski   |     70 | false          |
| tt0496406  | L'étrangère                    | Philippe Morier-Genoud |     70 | false          |
| tt0338330  | La passion de Bernadette       | Emmanuelle Riva        |     70 | false          |
| tt0095606  | The Music Teacher              | Philippe Volter        |     70 | false          |
| tt0184115  | Le petit voleur                | Yann Trégouët          |     70 | false          |
| tt4691166  | 7 Letters                      | Juliette Binoche       |     70 | false          |
| tt0115070  | ...à la campagne               | Benoît Régent          |     70 | false          |
| tt0109602  | Des feux mal éteints           | Hélène Vincent         |     70 | false          |
| tt0143022  | L'école est finie              | Hélène Vincent         |     70 | false          |
| tt0133063  | Lightmaker                     | Zbigniew Zamachowski   |     70 | false          |
| tt0113362  | The Horseman on the Roof       | Juliette Binoche       |     70 | true           |
| tt2474438  | Attila Marcel                  | Hélène Vincent         |     70 | true           |
| tt0083329  | Wielka majówka                 | Zbigniew Zamachowski   |     70 | false          |
| tt0436441  | Itinéraires                    | Yann Trégouët          |     70 | false          |
| tt2353767  | 1,000 Times Good Night         | Juliette Binoche       |     70 | true           |
| tt0066089  | The Modification               | Emmanuelle Riva        |     70 | false          |
| tt0054236  | Recourse in Grace              | Emmanuelle Riva        |     70 | false          |
| tt0093557  | Mon bel amour, ma déchirure    | Philippe Manesse       |     70 | false          |
| tt0094709  | The Gang of Four               | Benoît Régent          |     70 | false          |
| tt0223322  | Bitter Fruit                   | Emmanuelle Riva        |     70 | false          |
| tt0163776  | Egy tél az Isten háta mögött   | Florence Pernel        |     70 | false          |
| tt0056581  | Therese                        | Emmanuelle Riva        |     71 | false          |
| tt0191636  | The Widow of Saint-Pierre      | Juliette Binoche       |     71 | true           |
| tt0115658  | Bernie                         | Hélène Vincent         |     71 | true           |
| tt0241303  | Chocolat                       | Juliette Binoche       |     71 | true           |
| tt0103150  | Escape from the 'Liberty'      | Zbigniew Zamachowski   |     71 | false          |
|            | Cinema                         |                        |        |                |
| tt0412019  | Broken Flowers                 | Julie Delpy            |     71 | true           |
| tt0142383  | The Eighth Day                 | Emmanuelle Riva        |     71 | false          |
| tt0401711  | Paris, je t'aime               | Juliette Binoche       |     71 | true           |
| tt0836700  | Summer Hours                   | Isabelle Sadoyan       |     71 | true           |
| tt0104749  | Far from Brazil                | Emmanuelle Riva        |     71 | false          |
| tt0216625  | Code Unknown                   | Juliette Binoche       |     71 | true           |
| tt0104008  | A Tale of Winter               | Charlotte Véry         |     73 | true           |
| tt2209300  | Aftermath                      | Zbigniew Zamachowski   |     73 | true           |
| tt0198854  | Night of Destiny               | Philippe Volter        |     73 | false          |
| tt0091497  | Mauvais Sang                   | Juliette Binoche       |     73 | true           |
| tt0387898  | Caché                          | Juliette Binoche       |     73 | true           |
| tt0069145  | Somewhere, Someone             | Hugues Quester         |     73 | false          |
| tt0100878  | Una vita scellerata            | Florence Pernel        |     73 | false          |
| tt0084423  | La nuit de Varennes            | Hugues Quester         |     73 | true           |
| tt1020773  | Certified Copy                 | Juliette Binoche       |     73 | true           |
| tt0107260  | Joan the Maid 2: The Prisons   | Alain Ollivier         |     73 | false          |
| tt0096332  | The Unbearable Lightness of    | Juliette Binoche       |     73 | true           |
|            | Being                          |                        |        |                |
| tt0461254  | J'ai besoin d'air              | Yann Trégouët          |     73 | false          |
| tt0097106  | A Tale of Springtime           | Hugues Quester         |     73 | true           |
| tt0100263  | La Femme Nikita                | Jacques Disses         |     73 | true           |
| tt0086546  | City of Pirates                | Hugues Quester         |     73 | false          |
| tt0102421  | Mother                         | Isabelle Sadoyan       |     74 | true           |
| tt0090563  | Betty Blue                     | Claude Duneton         |     74 | true           |
| tt0113828  | Les Misérables                 | Isabelle Sadoyan       |     74 | true           |
| tt0262700  | Posseteni ot gospoda           | Philippe Volter        |     74 | false          |
| tt0116209  | The English Patient            | Juliette Binoche       |     74 | true           |
| tt0101318  | The Lovers on the Bridge       | Juliette Binoche       |     75 | true           |
| tt0052961  | Kapò                           | Emmanuelle Riva        |     75 | true           |
| tt0099334  | Cyrano de Bergerac             | Philippe Volter        |     75 | true           |
| tt0082949  | The Professional               | Pierre Forget          |     75 | true           |
| tt0119590  | Ma Vie en Rose                 | Hélène Vincent         |     75 | true           |
| tt8655470  | The Specials                   | Hélène Vincent         |     75 | true           |
| tt0053570  | Adua e le compagne             | Emmanuelle Riva        |     75 | false          |
| tt0111507  | Three Colors: White            | Julie Delpy            |     76 | true           |
| tt0119038  | Le Dîner de Cons               | Daniel Martin          |     76 | true           |
| tt0055082  | Léon Morin, Priest             | Emmanuelle Riva        |     76 | true           |
| tt0096163  | The Vanishing                  | Pierre Forget          |     76 | true           |
| tt0101765  | The Double Life of Véronique   | Claude Duneton         |     78 | true           |
| tt2209418  | Before Midnight                | Julie Delpy            |     79 | true           |
| tt1602620  | Amour                          | Emmanuelle Riva        |     79 | true           |
| tt0079322  | I... For Icarus                | Alain Ollivier         |     79 | true           |
| tt0052893  | Hiroshima Mon Amour            | Emmanuelle Riva        |     79 | true           |
| tt0108394  | Three Colors: Blue             | Alain Ollivier         |     79 | true           |
| tt0381681  | Before Sunset                  | Julie Delpy            |     80 | true           |
| tt0092593  | Au Revoir les Enfants          | Philippe Morier-Genoud |     80 | true           |
| tt0112471  | Before Sunrise                 | Julie Delpy            |     81 | true           |
| tt0111495  | Three Colors: Red              | Benoît Régent          |     81 | true           |
| tt0109682  | Du fond du coeur               | Benoît Régent          |     90 | false          |
+------------+--------------------------------+------------------------+--------+----------------+
>
```

# Value over Replacement Actor

It's interesting to ask how much an actor typically lifts or lowers a rating on IMDB. To answer this question, every actor's average rating of "Broadly Rated" (>1000 ratings) films on IMDB is taken, and stored. This is their "Rating". 

Then, for each film the actor appears in, the average rating of the all *other* actors in the film is taken, and the actual rating is compared. In this way, the actor either brings up or down the 'expected' rating of the film. 

You could do this for writers, directors, everyone involved in a film's production, but we don't do that here, just actors. 

To query the value system by range use ```v: [start index] [stop index]```.

This query is fundamentally subjective as to what cutoffs to make: we choose to require an actor appear in five Broadly Rated films to assess VOR. We ignore television, and anyone that does not have "actor" or "actress" as their first listed occupation in IMDB. 

The top few actors as of this writing show some interesting stories: Julie Dreyfus has only appeared in a few Tarantino films, and they are well rated on IMDB. We calculate her acting could lift your film's IMDB rating by 0.914 points; this could be significant! Because we don't add a director's average to the replacement calculations, Tarantino's films will tend to boost actors who only work with him. 

Long running "career" actors will tend to have lower VOR, as a function of being in so many movies. 

John Cazale is a great result, and shows that this system has some merit. He is best known for playing Fredo in The Godfather series, and is widely considered to have only taken excellent roles in excellent films. 

Edward Binns shows another quirk of our rules - his "famous films" include North by Northwest, 12 Angry Men, and others. His less famous films (and therefore rated less than 1000 times) include relative stinkers like "Oliver's Story" (rated 4.4). It doesn't have enough votes to bring down his VOR. 

Likewise, Abe Vigoda luckily avoids the penalty that comes from appearing in "Vasectomy: A Delicate Matter" (3.8). It's so bad, only a few people have seen it.

Two angles to go at for next round here; first would be to have different vote cutoffs for different decades; fewer people have seen films from the 50s than recent blockbusters. Alternately we could use some statistics; put another way, given 72 votes and a 3.8 rating, what is the probability that we have an indicative rating? Pretty high, I think.

```bash

> v:0 20
2020/05/13 11:20:58 Doing first fill of sorted Vor list
2020/05/13 11:20:58 Have 8861 Actors in the Vor
+-----------+----------------------+-------+
| nm0237838 | Julie Dreyfus        | 0.914 |
| nm0001030 | John Cazale          | 0.778 |
| nm0083081 | Edward Binns         | 0.711 |
| nm0100889 | Michael Bowen        | 0.688 |
| nm0893142 | Venkatesh Daggubati  | 0.687 |
| nm0015147 | Sitki Akçatepe       | 0.683 |
| nm0451425 | Kulbhushan Kharbanda | 0.640 |
| nm0001652 | John Ratzenberger    | 0.631 |
| nm0417301 | Janagaraj            | 0.619 |
| nm0109785 | Brent Briscoe        | 0.618 |
| nm0169454 | J.J. Cohen           | 0.612 |
| nm0001190 | David Prowse         | 0.608 |
| nm0001820 | Abe Vigoda           | 0.605 |
| nm0707399 | Rajendra Prasad      | 0.587 |
| nm1962192 | Upendra              | 0.580 |
| nm0017491 | Norman Alden         | 0.579 |
| nm0293285 | Alfonso Freeman      | 0.575 |
| nm1004985 | Yashpal Sharma       | 0.570 |
| nm0125084 | Paul Butler          | 0.565 |
| nm0411964 | Zeljko Ivanek        | 0.564 |
+-----------+----------------------+-------+
```

You can also look for actors that just bring things down when they appear:

```bash
> v:8840 8861
+-----------+----------------------+--------+
| nm2776304 | Ram Charan           | -0.673 |
| nm0001182 | Carmen Electra       | -0.675 |
| nm1139343 | Katt Williams        | -0.681 |
| nm0331374 | Marjoe Gortner       | -0.682 |
| nm0888727 | Musetta Vander       | -0.686 |
| nm0603352 | Dino Morea           | -0.693 |
| nm0891275 | Emmanuelle Vaugier   | -0.704 |
| nm1293381 | Sunny Leone          | -0.715 |
| nm0000137 | Bo Derek             | -0.758 |
| nm0641509 | Miles O'Keeffe       | -0.776 |
| nm0266422 | Jimmy Fallon         | -0.779 |
| nm1516058 | Deep Raj Rana        | -0.786 |
| nm1201734 | Peker Açikalin       | -0.791 |
| nm2365760 | Nikhil Dwivedi       | -0.817 |
| nm1863784 | Alp Kirsan           | -0.827 |
| nm0002246 | Harry Van Gorkum     | -0.849 |
| nm1300301 | Zayed Khan           | -0.881 |
| nm2596365 | Jacqueline Fernandez | -0.891 |
| nm5899377 | Tiger Shroff         | -0.960 |
| nm0370886 | Allison Hayes        | -0.987 |
+-----------+----------------------+--------+
```

I'm so sorry team.
