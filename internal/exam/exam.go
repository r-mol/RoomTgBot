package exam

import (
	"RoomTgBot/internal/consts"
	"RoomTgBot/internal/menus"
	"RoomTgBot/internal/mongodb"
	"RoomTgBot/internal/state"
	"RoomTgBot/internal/types"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"

	"context"

	"github.com/go-redis/redis/v9"
	telegram "gopkg.in/telebot.v3"
)

type Subject string

const (
	Subject0  Subject = ""
	Subject1          = "Compilers Construction"
	Subject2          = "Control Theory"
	Subject3          = "Data Mining"
	Subject4          = "Data Modeling and Databases I"
	Subject5          = "Data Modeling and Databases II"
	Subject6          = "Differential Equations"
	Subject7          = "Digital Signal Processing"
	Subject8          = "Distributed Systems"
	Subject9          = "Fundamentals of Computer Security"
	Subject10         = "Fundamentals of Robotics"
	Subject11         = "Game Theory"
	Subject12         = "Information Retrieval"
	Subject13         = "Information Theory"
	Subject14         = "Introduction to AI"
	Subject15         = "Introduction to Big Data"
	Subject16         = "Introduction to Machine Learning"
	Subject17         = "Lean Software Development"
	Subject18         = "Mechanics and Machines"
	Subject19         = "Networks"
	Subject20         = "Network and Cyber Security"
	Subject21         = "Non-Linear Optimization"
	Subject22         = "Operating Systems"
	Subject23         = "Philosophy II"
	Subject24         = "Physics I"
	Subject25         = "Practicum Project"
	Subject26         = "Probability and Statistics"
	Subject27         = "Robotics Systems"
	Subject28         = "Sensors and Sensing"
	Subject29         = "Software Architecture"
	Subject30         = "Software Project"
	Subject31         = "Software Systems Design"
	Subject32         = "System and Network Administration"
	Subject33         = "Theoretical Computer Science"
	Subject34         = "Theoretical Mechanics"
)

func (s Subject) ToId() [12]byte {
	switch s {
	case Subject1:
		return [12]byte{0}
	case Subject2:
		return [12]byte{1}
	case Subject3:
		return [12]byte{2}
	case Subject4:
		return [12]byte{3}
	case Subject5:
		return [12]byte{4}
	case Subject6:
		return [12]byte{5}
	case Subject7:
		return [12]byte{6}
	case Subject8:
		return [12]byte{7}
	case Subject9:
		return [12]byte{8}
	case Subject10:
		return [12]byte{9}
	case Subject11:
		return [12]byte{10}
	case Subject12:
		return [12]byte{11}
	case Subject13:
		return [12]byte{12}
	case Subject14:
		return [12]byte{13}
	case Subject15:
		return [12]byte{14}
	case Subject16:
		return [12]byte{15}
	case Subject17:
		return [12]byte{16}
	case Subject18:
		return [12]byte{17}
	case Subject19:
		return [12]byte{18}
	case Subject20:
		return [12]byte{19}
	case Subject21:
		return [12]byte{20}
	case Subject22:
		return [12]byte{21}
	case Subject23:
		return [12]byte{22}
	case Subject24:
		return [12]byte{23}
	case Subject25:
		return [12]byte{24}
	case Subject26:
		return [12]byte{25}
	case Subject27:
		return [12]byte{26}
	case Subject28:
		return [12]byte{27}
	case Subject29:
		return [12]byte{28}
	case Subject30:
		return [12]byte{29}
	case Subject31:
		return [12]byte{30}
	case Subject32:
		return [12]byte{31}
	case Subject33:
		return [12]byte{32}
	case Subject34:
		return [12]byte{33}
	}
	return [12]byte{}
}

func ToSubject(subject string) Subject {
	switch subject {
	case string(Subject1):
		return Subject1
	case string(Subject2):
		return Subject2
	case string(Subject3):
		return Subject3
	case string(Subject4):
		return Subject4
	case string(Subject5):
		return Subject5
	case string(Subject6):
		return Subject6
	case string(Subject7):
		return Subject7
	case string(Subject8):
		return Subject8
	case string(Subject9):
		return Subject9
	case string(Subject10):
		return Subject10
	case string(Subject11):
		return Subject11
	case string(Subject12):
		return Subject12
	case string(Subject13):
		return Subject13
	case string(Subject14):
		return Subject14
	case string(Subject15):
		return Subject15
	case string(Subject16):
		return Subject16
	case string(Subject17):
		return Subject17
	case string(Subject18):
		return Subject18
	case string(Subject19):
		return Subject19
	case string(Subject20):
		return Subject20
	case string(Subject21):
		return Subject21
	case string(Subject22):
		return Subject22
	case string(Subject23):
		return Subject23
	case string(Subject24):
		return Subject24
	case string(Subject25):
		return Subject25
	case string(Subject26):
		return Subject26
	case string(Subject27):
		return Subject27
	case string(Subject28):
		return Subject28
	case string(Subject29):
		return Subject29
	case string(Subject30):
		return Subject30
	case string(Subject31):
		return Subject31
	case string(Subject32):
		return Subject32
	case string(Subject33):
		return Subject33
	case string(Subject34):
		return Subject34
	}
	return Subject0
}

func GetSetExam(contex context.Context, bot *telegram.Bot, mdb *mongo.Client, rdb *redis.Client, ctx telegram.Context, subjectName string) error {
	curState, err := state.GetCurStateFromRDB(contex, rdb, ctx.Sender().ID)
	if err == redis.Nil {
		return ctx.Send("Please restart bot ✨")
	} else if err != nil {
		return err
	}

	if curState.StateName == consts.CommandExamDone {
		files := curState.Files
		photos := curState.Photos

		err = SetExam(contex, mdb, subjectName, files, photos)
		if err != nil {
			return err
		}

		curState.RemoveAll()

		commandFrom := curState.StateName

		err = state.CheckOfUserState(contex, rdb, ctx, commandFrom, consts.CommandStart)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Exams successful set...", menus.MainMenu)
	} else if curState.StateName == consts.CommandGetExam {
		curState.Files, curState.Photos, err = GetExam(contex, mdb, subjectName)
		if err != nil {
			return err
		}

		err = curState.SendAllAvailableMessage(bot, ctx.Sender(), state.Message{}, menus.MainMenu)
		if err != nil {
			return err
		}

		curState.RemoveAll()

		commandFrom := curState.StateName

		err = state.CheckOfUserState(contex, rdb, ctx, commandFrom, consts.CommandStart)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Exams successful get...", menus.MainMenu)
	}

	return ctx.Send("Please restart bot ✨")
}

func getExam(ctx context.Context, client *mongo.Client, subjectName string) (types.ExamEntry, error) {
	mongoExams, err := mongodb.GetAll[types.ExamEntry](ctx, client, consts.MongoExamCollection)
	if err != nil {
		return types.ExamEntry{}, fmt.Errorf("unable to get Exams from mongodb: %v", err)
	}

	exam := types.ExamEntry{}

	for _, tmpExam := range mongoExams {
		if tmpExam.MongoID == ToSubject(subjectName).ToId() {
			exam = tmpExam
			break
		}
	}

	return exam, nil
}

func setExam(ctx context.Context, client *mongo.Client, subjectName string, files []telegram.Document, photos []telegram.Photo) error {
	mongoExams, err := mongodb.GetAll[types.ExamEntry](ctx, client, consts.MongoExamCollection)
	if err != nil {
		return fmt.Errorf("unable to get Exams from mongodb: %v", err)
	}

	exam := types.ExamEntry{}
	flag := false

	for _, tmpExam := range mongoExams {
		if tmpExam.MongoID == ToSubject(subjectName).ToId() {
			exam = tmpExam
			flag = true
			break
		}
	}

	if flag {
		exam.MongoID = ToSubject(subjectName).ToId()
	}

	exam.Files.Files = append(exam.Files.Files, files...)
	exam.Files.Photos = append(exam.Files.Photos, photos...)

	_, err = mongodb.AddOne(ctx, client, consts.MongoExamCollection, &exam)

	if err != nil {
		return fmt.Errorf("unable to add exam to mongodb: %v", err)
	}
	return nil
}

func SetExam(ctx context.Context, client *mongo.Client, subjectName string, files []telegram.Document, photos []telegram.Photo) error {
	return setExam(ctx, client, subjectName, files, photos)
}

func GetExam(ctx context.Context, client *mongo.Client, subjectName string) ([]telegram.Document, []telegram.Photo, error) {
	exam, err := getExam(ctx, client, subjectName)
	if err != nil {
		return []telegram.Document{}, []telegram.Photo{}, err
	}
	return exam.Files.Files, exam.Files.Photos, nil
}
