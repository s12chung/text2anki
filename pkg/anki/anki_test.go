package anki

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/ioutil"
	"github.com/s12chung/text2anki/pkg/util/logg"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
	"github.com/s12chung/text2anki/pkg/util/test/fixture/flog"
)

var plog = flog.FixtureUpdateNoWrite()

func init() {
	dir := path.Join(os.TempDir(), test.GenerateName("anki.TestMain"))
	c := Config{ExportPrefix: "t2a-", NotesCacheDir: dir}
	if err := os.MkdirAll(dir, ioutil.OwnerRWXGroupRX); err != nil {
		plog.Error("anki.init()", logg.Err(err)) //nolint:forbidigo // used in init only
		os.Exit(-1)
	}
	SetConfig(c)
}

func notesFromFixture(t *testing.T) []Note {
	require := require.New(t)

	bytes, err := os.ReadFile(fixture.JoinTestData("Notes.json"))
	require.NoError(err)
	var notes []Note
	require.NoError(json.Unmarshal(bytes, &notes))
	for _, note := range notes {
		test.EmptyFieldsMatch(t, note, "usageSoundSource")
	}
	return notes
}

func notesWithSounds(t *testing.T) []Note {
	require := require.New(t)

	notes := notesFromFixture(t)
	sound := fixture.Read(t, "sound.mp3")
	for i := range notes {
		if i%2 == 1 {
			continue
		}
		err := notes[i].SetSound(sound, fmt.Sprintf("Naver CLOVA Speech Synthesis - %v", i))
		require.NoError(err)
	}
	return notes
}

type soundFactory struct{}

func (s soundFactory) Name() string { return "soundFactory name" }
func (s soundFactory) Sound(_ context.Context, usage string) ([]byte, error) {
	return []byte(usage), nil
}

func TestSoundSetter_SetSound(t *testing.T) {
	require := require.New(t)

	soundSetter := NewSoundSetter(soundFactory{})
	notes := notesFromFixture(t)
	require.NoError(soundSetter.SetSound(context.Background(), notes))
	for _, note := range notes {
		require.True(note.HasUsageSound())
		require.Equal(soundSetter.soundFactory.Name(), note.usageSoundSource)
		require.Equal(note.Usage, string(test.Read(t, path.Join(config.NotesCacheDir, note.UsageSoundFilename()))))
	}
}

func TestNote_ID(t *testing.T) {
	require := require.New(t)
	require.Equal("어른-Flower Road-모자람 없이 주신 사랑이 과분하다 느낄 때쯤 난 어른이 됐죠", notesFromFixture(t)[0].ID())
}

func TestNote_SetSound(t *testing.T) {
	require := require.New(t)

	note := notesFromFixture(t)[0]
	require.False(note.HasUsageSound())

	soundSource := "the source"
	soundContents := []byte("my_test")
	require.NoError(note.SetSound(soundContents, soundSource))

	require.True(note.HasUsageSound())
	require.Equal(soundSource, note.usageSoundSource)
	require.Equal(soundContents, test.Read(t, path.Join(config.NotesCacheDir, note.UsageSoundFilename())))
}

func TestNote_CSV(t *testing.T) {
	testName := "TestNote_CSV"
	fixture.CompareReadOrUpdate(t, testName+".json", test.JSON(t, notesFromFixture(t)[0].CSV()))
}

func TestNote_UsageSoundFilename(t *testing.T) {
	tcs := []struct {
		usage  string
		result string
	}{
		{usage: "여러 가지 야채들, 그리고 달걀, 심지어는 불고기, 김치 등등 다양한 음식을 넣기 시작했다고 합니다.", result: "t2a-여러 가지 야채들, 그리고 달걀, 심지어는 불고기, 김치 등등 다양한 음식을 넣기 시.mp3"},
		{usage: "“김밥은 다양한 야채가 골고루 들어가 있고, 고기도 들어가 있으니, 정말 건강에 좋은 거니까", result: "t2a-“김밥은 다양한 야채가 골고루 들어가 있고, 고기도 들어가 있으니, 정말 건강에 .mp3"},
		{usage: "그래서 영양분도 많고, 피도 맑게 해 주는 그런 미역을 국으로 끓여서 먹기 시작한 거죠.", result: "t2a-그래서 영양분도 많고, 피도 맑게 해 주는 그런 미역을 국으로 끓여서 먹기 시작한 .mp3"},
		{usage: "그 이유는 엄마가 나를 낳았다는 사실을 항상 감사하라는 의미에서 미역국을 먹기 시작했다고 합니다.", result: "t2a-그 이유는 엄마가 나를 낳았다는 사실을 항상 감사하라는 의미에서 미역국을 먹기 .mp3"},
		{usage: "이렇게 드라마가 유행을 하게 되면 그 드라마에 나온 스타일을 다 따라 하는 것 같아요.", result: "t2a-이렇게 드라마가 유행을 하게 되면 그 드라마에 나온 스타일을 다 따라 하는 것 같.mp3"},
		{usage: "‘아! 송혜교가 발랐었지. 그때 굉장히 예뻤는데.’ 하면서 저도 모르게 이렇게 집고 있더라고요.", result: "t2a-‘아! 송혜교가 발랐었지. 그때 굉장히 예뻤는데.’ 하면서 저도 모르게 이렇게 집.mp3"},
		{usage: "어때요? 어울리나요?", result: "t2a-어때요 어울리나요.mp3"},
		{usage: "그렇지만 대부분의 한국 사람들이 유행을 따라서 옷을 입고, 화장을 하고, 악세사리를 하는 것 같습니다.", result: "t2a-그렇지만 대부분의 한국 사람들이 유행을 따라서 옷을 입고, 화장을 하고, 악세사리.mp3"},
		{usage: "그런데 옷 가게나 화장품 가게의 주인들이 엄청 빠르게 그런 스타일의 옷을 갖다 놔요.", result: "t2a-그런데 옷 가게나 화장품 가게의 주인들이 엄청 빠르게 그런 스타일의 옷을 갖다 놔.mp3"},
		{usage: "맞아요. 그래도 결혼하기 전에 자기가 좋아하는 스타일을 아주 구체적으로 정하면 좋아요.", result: "t2a-맞아요. 그래도 결혼하기 전에 자기가 좋아하는 스타일을 아주 구체적으로 정하면 .mp3"},
		{usage: "진짜요? 그럼 헷갈리지 않아요?", result: "t2a-진짜요 그럼 헷갈리지 않아요.mp3"},
		{usage: "꿈을 안 꾸고 자면 좋을 것 같은데 항상 꿈을 꾸니까 저도 뭐 어쩔 수 없을 것 같아요.", result: "t2a-꿈을 안 꾸고 자면 좋을 것 같은데 항상 꿈을 꾸니까 저도 뭐 어쩔 수 없을 것 같아.mp3"},
		{usage: "아 그래요? 그럼 말 놓을까요?", result: "t2a-아 그래요 그럼 말 놓을까요.mp3"},
		{usage: "아 진짜요? 그 언니들은 먼저 말 놨어요?", result: "t2a-아 진짜요 그 언니들은 먼저 말 놨어요.mp3"},
		{usage: "어… 이렇게 공항에 오면 큰 가방을 가지고 이렇게 왔다 갔다 바쁘게 왔다 갔다 하는 사람들을 보면", result: "t2a-어… 이렇게 공항에 오면 큰 가방을 가지고 이렇게 왔다 갔다 바쁘게 왔다 갔다 하.mp3"},
		{usage: "아, 그래요? 결혼하라는 압박 그런 건 없구요?", result: "t2a-아, 그래요 결혼하라는 압박 그런 건 없구요.mp3"},
		{usage: "어. 그래서 엘에이에 가 가지고 맛있는 거 많이 먹고 사람들이 날 알아봐 주고 막 그랬던 게 막 생각나네.", result: "t2a-어. 그래서 엘에이에 가 가지고 맛있는 거 많이 먹고 사람들이 날 알아봐 주고 막 그.mp3"},
		{usage: "이렇게 말을 할 때 숨소리가 너무 마이크에 심하게 들어가지 않도록 해 주는 필터가 있어요.", result: "t2a-이렇게 말을 할 때 숨소리가 너무 마이크에 심하게 들어가지 않도록 해 주는 필터가.mp3"},
		{usage: "현재는 따로 다니는 직장은 없고, 프리랜서로 그림 그리는 일 하고 있어요. 지연 씨는요?", result: "t2a-현재는 따로 다니는 직장은 없고, 프리랜서로 그림 그리는 일 하고 있어요. 지연 씨.mp3"},
		{usage: "아, 그러셨어요? 식사는 하셨나요?", result: "t2a-아, 그러셨어요 식사는 하셨나요.mp3"},
		{usage: "그럼 따뜻한 차라도 한 잔 드릴까요? 녹차? 커피?", result: "t2a-그럼 따뜻한 차라도 한 잔 드릴까요 녹차 커피.mp3"},
		{usage: "요즘 왜 이렇게 조용해? 잘 지내?", result: "t2a-요즘 왜 이렇게 조용해 잘 지내.mp3"},
		{usage: "남자 친구 주말에 쉬어? 원래 평일에 쉬지 않아?", result: "t2a-남자 친구 주말에 쉬어 원래 평일에 쉬지 않아.mp3"},
		{usage: "목걸이? 너 작년에도 남자 친구한테 목걸이 받지 않았어?", result: "t2a-목걸이 너 작년에도 남자 친구한테 목걸이 받지 않았어.mp3"},
		{usage: "우리도 아무 계획이 없어서, 신랑이랑 얘기했는데, 친한 사람들 초대 해서 우리 집에서 간단하게 식사하면서 술 마시면 어떨까 해.", result: "t2a-우리도 아무 계획이 없어서, 신랑이랑 얘기했는데, 친한 사람들 초대 해서 우리 집.mp3"},
		{usage: "그래. 은경이도 좋아할 거야. 우리는 뭐 사갈까? 필요한 거 없어?", result: "t2a-그래. 은경이도 좋아할 거야. 우리는 뭐 사갈까 필요한 거 없어.mp3"},
		{usage: "너네 둘이 싸워서? 그거 완전 옛날 일 아니야?", result: "t2a-너네 둘이 싸워서 그거 완전 옛날 일 아니야.mp3"},
		{usage: "절대 데려오면 안 된다! 알았지? 어?", result: "t2a-절대 데려오면 안 된다! 알았지 어.mp3"},
		{usage: "아인아, 오늘 몇 시에 들어와? 오늘도 늦어?", result: "t2a-아인아, 오늘 몇 시에 들어와 오늘도 늦어.mp3"},
		{usage: "네? 진짜요? 회사에서 저녁 먹고 들어오면 안 돼요?", result: "t2a-네 진짜요 회사에서 저녁 먹고 들어오면 안 돼요.mp3"},
		{usage: "뭐? 너 시험이 내일모레인데 왜 늦어? 오늘 늦기만 해 봐.", result: "t2a-뭐 너 시험이 내일모레인데 왜 늦어 오늘 늦기만 해 봐..mp3"},
		{usage: "안 해. 저번에도 잘생겼다고 해서 나갔더니 완전 폭탄이었으면서. 엄마 말 이제 안 믿어.", result: "t2a-안 해. 저번에도 잘생겼다고 해서 나갔더니 완전 폭탄이었으면서. 엄마 말 이제 안 .mp3"},
		{usage: "엄마, 나 나가요. 나 요즘 회사에서 피곤해서 집에 와서도 아무것도 못하고 잠만 자잖아.", result: "t2a-엄마, 나 나가요. 나 요즘 회사에서 피곤해서 집에 와서도 아무것도 못하고 잠만 자.mp3"},
		{usage: "여섯 시? 아침 여섯 시?", result: "t2a-여섯 시 아침 여섯 시.mp3"},
		{usage: "진짜요? 아, 근데 뭔가 이상한데?", result: "t2a-진짜요 아, 근데 뭔가 이상한데.mp3"},
		{usage: "아, 그래요? 고민되네, 어쩌지?", result: "t2a-아, 그래요 고민되네, 어쩌지.mp3"},
		{usage: "그래? 작은 거 같아? 좀 이상한 거 같긴 해.", result: "t2a-그래 작은 거 같아 좀 이상한 거 같긴 해..mp3"},
		// this is a case that I can't figure out...
		// {usage: "네?? 하... 정말... 같은 수업을 옆자리 앉아서 몇 개월을 같이 들었는데 어떻게 이름도 몰라요?", result: "t2a-네 하... 정말... 같은 수업을 옆자리 앉아서 몇 개월을 같이 들었는데 어떻게 이름.mp3"},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.usage, func(t *testing.T) {
			require := require.New(t)
			note := Note{Usage: tc.usage}
			require.Equal(tc.result, note.UsageSoundFilename())
		})
	}
}

func TestNote_HasUsageSound(t *testing.T) {
	require := require.New(t)

	note := notesFromFixture(t)[0]
	require.False(note.HasUsageSound())

	soundSource := "the source"
	soundContents := []byte("my_test")
	require.NoError(note.SetSound(soundContents, soundSource))

	require.True(note.HasUsageSound())
}

func TestExportFiles(t *testing.T) {
	require := require.New(t)
	testName := "TestExportFiles"

	exportDir := path.Join(os.TempDir(), test.GenerateName(testName))
	test.MkdirAll(t, exportDir)
	require.NoError(ExportFiles(exportDir, notesWithSounds(t)))

	fixture.CompareReadOrUpdateDir(t, "ExportFiles", exportDir)
}

func TestExportSounds(t *testing.T) {
	require := require.New(t)
	testName := "TestExportSounds"

	exportDir := path.Join(os.TempDir(), test.GenerateName(testName))
	test.MkdirAll(t, exportDir)
	require.NoError(ExportSounds(exportDir, notesWithSounds(t)))

	dirEntries, err := os.ReadDir(exportDir)
	require.NoError(err)
	dirEntryNames := make([]string, len(dirEntries))
	for i, dirEntry := range dirEntries {
		dirEntryNames[i] = dirEntry.Name()
	}
	require.Equal([]string{"t2a-꽃길만 걷게 해줄게요.mp3", "t2a-모자람 없이 주신 사랑이 과분하다 느낄 때쯤 난 어른이 됐죠.mp3"}, dirEntryNames)
	for _, entry := range dirEntries {
		require.Equal("sound.mp3 fake", string(test.Read(t, path.Join(exportDir, entry.Name()))))
	}
}

func TestExportCSVFile(t *testing.T) {
	require := require.New(t)
	testName := "TestExportCSVFile"

	dir := path.Join(os.TempDir(), test.GenerateName(testName))
	test.MkdirAll(t, dir)

	dir = path.Join(dir, "basic.csv")
	require.NoError(ExportCSVFile(dir, notesFromFixture(t)))
	fixture.CompareReadOrUpdate(t, testName+".csv", test.Read(t, dir))
}

func TestExportCSV(t *testing.T) {
	require := require.New(t)

	buffer := &bytes.Buffer{}
	err := ExportCSV(buffer, notesFromFixture(t))
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, "TestExportCSVFile.csv", buffer.Bytes())
}
