package route

// cspell:ignore curri, subscr, signup

import (
	"dazzlefm/service"
	"dazzlefm/service/form"
	"dazzlefm/util"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

// Schedule ...
type Schedule struct {
	log *zap.SugaredLogger
	env *util.Environment

	Office service.IOffice
}

const longForm = "2006-01-02 15:04:05 -0700 MST"

// Init ...
func (s *Schedule) Init(prefix string, env *util.Environment) error {
	s.env = env
	s.log = env.Log.Sugar().Named("schedule")
	s.Office = env.Service("office").(service.IOffice)

	nPre := path.Join(prefix, "schedule")
	grp := env.Rtr.Group(nPre)

	grp.GET("", s.GetList)
	grp.GET("/count", s.GetListCount)
	grp.GET("/:id", s.Get)
	grp.GET("/", s.GetAList)
	grp.POST("", s.Add)
	grp.POST("/:id", s.Save)
	grp.DELETE("/:id", s.Delete)

	grp.GET("/summary", s.ScheduleSummary)
	grp.GET("/all", s.GetDay)

	return nil
}

// GetList GET: /schedule
func (s Schedule) GetList(c echo.Context) (err error) {

	// if ok, err := IsAuthorized(c, 1, ScheduleModule); err != nil || !ok {
	// 	if err != nil {
	// 		s.log.Error(err)
	// 		return err
	// 	}
	// 	return c.String(http.StatusNotFound, "!OK")
	// }

	page := c.QueryParam("page")

	pageNo, err := strconv.Atoi(page)
	if err != nil {
		return
	}

	offSet := c.QueryParam("offset")

	offset, err := strconv.Atoi(offSet)
	if err != nil {
		return
	}

	dayID := c.QueryParam("day")

	day, err := strconv.Atoi(dayID)
	if err != nil {
		return err
	}

	resp := util.Response{Code: 200}
	list, err := s.Office.ListSchedule(pageNo, offset, day)
	if err != nil {
		return
	}

	resp.Set("list", list)
	return c.JSON(200, resp)
}

// GetListCount GET: /schedule/count
func (s Schedule) GetListCount(c echo.Context) (err error) {

	resp := util.Response{Code: 200}
	cnt, err := s.Office.CountSchedule()
	if err != nil {
		return
	}

	resp.Set("count", cnt)
	return c.JSON(200, resp)
}

// GetAList GET: /schedule/?presenterid=id
func (s Schedule) GetAList(c echo.Context) (err error) {

	// if ok, err := IsAuthorized(c, 1, ScheduleModule); err != nil || !ok {
	// 	if err != nil {
	// 		s.log.Error(err)
	// 		return err
	// 	}
	// 	return c.String(http.StatusNotFound, "!OK")
	// }

	presenterID := c.QueryParam("presenterid")
	// log.Println(presenterID)

	pID, err := strconv.ParseInt(presenterID, 10, 64)

	if err != nil {
		log.Println(err)
		return err
	}
	// log.Println(pID)
	// pID := util.Int64Param(c, presenterID)

	resp := util.Response{Code: 200}
	list, err := s.Office.ListASchedule(pID)
	if err != nil {
		return
	}

	resp.Set("list", list)
	return c.JSON(200, resp)
}

// Get GET: /schedule/id
func (s Schedule) Get(c echo.Context) (err error) {

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	record, err := s.Office.GetSchedule(id)
	if err != nil {
		return
	}

	resp.Set("schedule", record)
	return c.JSON(200, resp)
}

// Add POST: /schedule
func (s Schedule) Add(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 2, ScheduleModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{}

	form := &form.Schedule{}
	if err = c.Bind(form); err != nil {
		s.log.Error(err)
		return
	}

	s.log.Debugf("%v\n", form.StartTime)

	//timeN = fmt.Sprintf("1970-02-02T%sZ", timeNow)
	loc, _ := time.LoadLocation("Africa/Lagos")
	start, err := time.ParseInLocation(longForm, fmt.Sprintf("%v", form.StartTime), loc)
	if err != nil {
		s.log.Errorf("%v", err)
		return err
	}

	end, err := time.ParseInLocation(longForm, fmt.Sprintf("%v", form.EndTime), loc)
	if err != nil {
		s.log.Errorf("%v", err)
		return err
	}

	form.StartTime = start
	form.EndTime = end

	exist, res, err := s.Office.AddSchedule(form)
	if err != nil {
		return
	}

	if exist {
		resp.Code = http.StatusConflict
		resp.ErrMsg("exists", res.Title)
		return c.JSON(200, resp)
	}
	resp.Code = 200
	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Save POST: /schedule/id
func (s Schedule) Save(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 3, ScheduleModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{}
	id := util.Int64Param(c, "id")

	form := &form.Schedule{}
	if err = c.Bind(form); err != nil {
		return
	}

	loc, _ := time.LoadLocation("Africa/Lagos")
	start, err := time.ParseInLocation(longForm, fmt.Sprintf("%v", form.StartTime), loc)
	if err != nil {
		return err
	}

	end, err := time.ParseInLocation(longForm, fmt.Sprintf("%v", form.EndTime), loc)
	if err != nil {
		return err
	}

	form.StartTime = start
	form.EndTime = end

	exist, res, err := s.Office.SaveSchedule(id, form)
	if err != nil {
		return
	}

	if exist {
		s.log.Debugf("%s", res.Title)
		resp.Code = http.StatusConflict
		resp.ErrMsg("exists", res.Title)
		return c.JSON(200, resp)
	}
	resp.Code = 200
	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Delete DELETE: /schedule/id
func (s Schedule) Delete(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 4, ScheduleModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	err = s.Office.DeleteSchedule(id)
	if err != nil {
		return
	}

	resp.Set("id", id)
	return c.JSON(200, resp)
}

// ScheduleSummary GET: /schedule/summary
func (s Schedule) ScheduleSummary(c echo.Context) (err error) {

	// if ok, err := IsAuthorized(c, 1, ScheduleModule); err != nil || !ok {
	// 	if err != nil {
	// 		s.log.Error(err)
	// 		return err
	// 	}
	// 	return c.String(http.StatusNotFound, "!OK")
	// }

	resp := util.Response{Code: 200}

	day := c.QueryParam("day")
	dayValue, err := strconv.Atoi(day)
	if err != nil {
		return err
	}

	timeNow := c.QueryParam("time")
	timeWithTz := timeNow + "+1"
	schedules, err := s.Office.GetScheduleSummary(dayValue, timeWithTz)
	if err != nil {
		return
	}

	// timeNow = strings.Replace(timeNow, ":", "", -1)
	// Agreed date := 1970-02-02 which is equivalent to 020270
	// agreedDate := "020270"
	// Comparing it with
	// Mon Jan 2 15:04:05 -0700 MST 2006  or
	// 01/02 03:04:05PM '06 -0700
	timeNow = fmt.Sprintf("1970-02-02T%s+01:00", timeNow)
	s.log.Debugf("%v\n", timeNow)
	loc, err := time.LoadLocation("Africa/Lagos")
	if err != nil {
		s.log.Error(err)
		return err
	}
	newTime, err := time.ParseInLocation(time.RFC3339, timeNow, loc)
	if err != nil {
		s.log.Error(err)
		return err
	}

	log.Printf("newTime:%v\n", newTime)
	//log.Printf("schedules[0].StartTime:%v\n", schedules[0].StartTime)
	// log.Printf("schedules[1].StartTime:%v\n", schedules[1].StartTime)

	if len(schedules) > 0 {
		if newTime.Before(schedules[0].StartTime) {
			schedules[0].Caption = form.CTNext
		} else if schedules[0].StartTime.Before(newTime) {

			schedules[0].Caption = form.CTLiveNow
		}

		if schedules[0].Caption == form.CTLiveNow {
			if len(schedules) > 1 {
				schedules[1].Caption = form.CTNext
			}

		} else {
			if len(schedules) > 1 {
				schedules[1].Caption = form.CTOthers
			}

		}

		for i := range schedules {
			if i > 1 {
				schedules[i].Caption = form.CTOthers
			}

		}
	}

	resp.Set("others", schedules)
	return c.JSON(200, resp)
}

// GetDay GET: /all?day=1
func (s Schedule) GetDay(c echo.Context) (err error) {

	// if ok, err := IsAuthorized(c, 1, ScheduleModule); err != nil || !ok {
	// 	if err != nil {
	// 		s.log.Error(err)
	// 		return err
	// 	}
	// 	return c.String(http.StatusNotFound, "!OK")
	// }

	resp := util.Response{Code: 200}

	day := c.QueryParam("day")
	dayValue, err := strconv.Atoi(day)
	if err != nil {
		return err
	}

	schedules, err := s.Office.GetAllSchedules(dayValue)
	if err != nil {
		return
	}

	resp.Set("all", schedules)
	return c.JSON(200, resp)
}
