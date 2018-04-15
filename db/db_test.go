package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func InsertTestData() {
	// Insert users
	AddUser("alice123", "alicepassword", "127.0.0.1")
	AddUser("bob", "Password", "123.456.789")
	AddUser("karateAMD", "pwd123", "192.168.10.123")
	AddUser("sameetandpotatoes", "iLuvMacs", "10.192.345.987")
	AddUser("archilmedes", "linuxFTW", "987.654.321")
	AddUser("andrew", "anotherPass", "888.888.888")

	//Insert sessions
	InsertIntoSessions(12, 1, 2, "3478462413678237ab87846754785489329853e47237646718487980423f095874236784675889490543874123675478329056bc7823619560458372e956", "12D345678902F83AE")
	InsertIntoSessions(14, 1, 4, "458747cb3457654765daf687536899857846734674735896590efbd8564328759203185487398cd39574293143478375a485695864585907568458857438", "CBDEABD347ABDC392")
	InsertIntoSessions(35, 3, 5, "abcdb378675934bdbd0935847349036985fbd490590584374374b43894784578431243b37465723894d3434981fdcb484726739923874bd3837473fedb31", "ACD537492B126CD43")
	InsertIntoSessions(34, 3, 4, "f4783467463b347458ab7687989e876c876778988767b87678f7986878a878e67bc34734763fb3274632cdba4348734384b34873362ab834728392847343", "FFBC63729D82635EC")
	InsertIntoSessions(64, 6, 4, "bcd8763728749378ab8347839847328492ae897638903478b834743898c834738423e9786f9ff7657843cb3874383b8973487ef3864727384a8783647873", "B3847C837D77654E5")
	InsertIntoSessions(62, 6, 2, "abc78384689234752369071625d8736543976d7f967798b567789098768789076890876890786890876890e8789087e877a90887b89c87d77e762532dbcf", "675C6A7CA877B6A67")
	InsertIntoSessions(32, 3, 2, "903873473785b4084787d8767889e988767543c45655434567a56467897669c009987766565b78765545e7896767f3234674734589bc9084734efb398473", "BC5A56DBDE836464B")

	//Insert conversations
	InsertIntoConversations(12, "Hello World", "01/02/2017:08:20:19", 1)
	InsertIntoConversations(14, "Hey Sameet, its Alice <3", "02/14/2018:11:11:11", 0)
	InsertIntoConversations(35, "Hey Andrew, I need help with 511, when are you free?", "04/10/2018:12:30:08", 1)
	InsertIntoConversations(52, "lul", "03/28/2018:18:04:10", 0)
	InsertIntoConversations(35, "I almost made my Mac a brick", "04/08/2018:17:01:40", 1)
	InsertIntoConversations(42, "Why did the chicken cross the road?", "04/12/2018:07:56:00", 1)
	InsertIntoConversations(42, "To get to the other side?", "04/12/2018:07:59:13", 0)
	InsertIntoConversations(35, "When are we playing Fortnite?", "04/08/2018:17:59:02", 0)

}
func SetupDatabaseTest(t *testing.T) {
	//DB := SetupDatabase()
	SetupDatabase()
	tables := ShowTables()
	assert.Equal(t, 3, len(tables))
	//return DB
}

func TestDatabase(t *testing.T) {
	SetupDatabaseTest(t)
	InsertTestData()
	UsersTest(t)
	ConversationTest(t)
	SessionsTest(t)
}
