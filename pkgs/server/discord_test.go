package server_test

// func TestRegisterCommands(t *testing.T) {
// 	//Create an in memory database
// 	db, err := server.NewDB(&config.Config{
// 		DBPath: ":memory:",
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = server.InitDB(db)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Add a few commands to the database
// 	commands := []server.CommandsTableRow{
// 		{
// 			Command: "test",
// 			Hash:    "test",
// 			Guildid: 170605932986761216,
// 		},
// 		{
// 			Command: "test2",
// 			Hash:    "test2",
// 			Guildid: 170605932986761216,
// 		},
// 		{
// 			Command: "test3",
// 			Hash:    "test3",
// 			Guildid: 2,
// 		},
// 	}

// 	// Mock a new discord session
// 	state, err := newState()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	session, err = mocksession.New(
// 		mocksession.WithState(state),
// 		mocksession.WithClient(&http.Client{
// 			Transport: mockrest.NewTransport(state),
// 		}),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Call RegisterCommands
// 	err = server.RegisterCommands(db, session)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// }
