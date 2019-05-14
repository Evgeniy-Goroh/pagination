package Paginator

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Paginator(t *testing.T) {
	Convey("Basic logics", t, func() {
		p := New(20)
		So(len(p.Pages()), ShouldEqual, 2)
	})

	Convey("Custom logics", t, func() {
		p := Custom(&Config{PageSize: 10, Current: 2, LinkedCount: 3}, 23)
		So(len(p.Pages()), ShouldEqual, 3)
		So(p.TotalPages(), ShouldEqual, 3)
		So(p.IsFirst(), ShouldBeFalse)
		So(p.HasPrevious(), ShouldBeTrue)
		So(p.Previous(), ShouldEqual, 1)
		So(p.HasNext(), ShouldBeTrue)
		So(p.Next(), ShouldEqual, 3)
		So(p.IsLast(), ShouldBeFalse)

		Convey("LinkedCount", func() {
			p := Custom(&Config{PageSize: 10, Current: 2, LinkedCount: 0}, 23)
			So(len(p.Pages()), ShouldEqual, 0)

			p = Custom(&Config{PageSize: 10, Current: 2, LinkedCount: 0}, 5)
			So(len(p.Pages()), ShouldEqual, 0)
			So(p.Current(), ShouldEqual, 1)
		})

		Convey("LinkedCount TotalPages", func() {
			p := Custom(&Config{PageSize: 10, Current: 2, LinkedCount: 1}, 23)
			So(len(p.Pages()), ShouldEqual, 1)

			p = Custom(&Config{PageSize: 10, Current: 2, LinkedCount: 1}, 5)
			So(len(p.Pages()), ShouldEqual, 1)
		})

		Convey("TotalPages LinkedCount", func() {
			Convey("LinkedCount test", func() {
				Print("\n")
				p := Custom(&Config{PageSize: 10, Current: 1, LinkedCount: 3}, 63)
				for _, page := range p.Pages() {
					Printf("%v ", page)
				}
				Print("\n")

				p = Custom(&Config{PageSize: 10, Current: 3, LinkedCount: 3}, 63)
				for _, page := range p.Pages() {
					Printf("%v ", page)
				}
				Print("\n")

				p = Custom(&Config{PageSize: 10, Current: 6, LinkedCount: 3}, 63)
				for _, page := range p.Pages() {
					Printf("%v ", page)
				}
				Print("\n")
			})

			Convey("LinkedCount test2", func() {
				Print("\n")
				p := Custom(&Config{PageSize: 10, Current: 1, LinkedCount: 4}, 63)
				for _, page := range p.Pages() {
					Printf("%v ", page)
				}
				Print("\n")

				p = Custom(&Config{PageSize: 10, Current: 4, LinkedCount: 4}, 63)
				for _, page := range p.Pages() {
					Printf("%v ", page)
				}
				Print("\n")

				p = Custom(&Config{PageSize: 10, Current: 6, LinkedCount: 4}, 63)
				for _, page := range p.Pages() {
					Printf("%v ", page)
				}
				Print("\n")
			})
		})
	})
}
