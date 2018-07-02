Simple Profiler for timing sections of go code

If you want to optimize code, it helps knowing where to start looking for improvements. This profiler allows you to 
mark bits of code (using a Tick(...) call), and if active will measure and report how long it takes each sub-segment to execute.
You can mark as many or as few lines of code as you like, as well as set the level at which that particular profiler will activate.
You can have multiple profilers running using different names, with ticks between two points you want to measure within a function.
It has been designed to be as painless to add as possible. 

Once you've setup a profiler in the code, you can activate it by setting a level at runtime (usually via command line arguments), and
all of the profilers at or below that level will be active.

Inside your function, when you are creating a new instance of the profiler you decide at what level the section of code should profile. 
In production I might use level 1 for a html GET requesthandler, level 2 inside important functions, and 3 for line by line breakdowns
of potentially problematic code. If I want to see how things are generally running, I'll run: myapp -p 1, which will show me the big
functions. If it seems a bit funny I'll run with -p 1 or 2 depending on how deep I want to profile.

You call the profiler.New function with 2 arguments; the name of the profiler (for display purposes) and the level. 0 or negative levels
will always show, a level 1 will show when the global level is 1 or higher. The idea is to create your profiling code and categorize it 
once, and when you want to watch the profiling data, you run the app with a -p1 or -p5 flag depending on how much detail you want
(only useful if you've set profilers up going to level 5).

It's important to note, that when a function runs, each profiler will still need to check whether it should record data, so it will 
definitely consume some unnecesssary cpu cycles. It does not record data unless it is active, but it does do a few comparison checks 
(if true, if mylevel < globallevel, that sort of thing). If you are after absolute speed, then after identifying and optimizing 
problematic code you may consider removing the profiling code as you see fit. 

Example Usage:

in package main:

import (
    "flag"
    "github.com/zeroepix/goprofiler"
)

var (
    profileLevelFlag     = flag.Int("p", 0, "choose the profiling level to log and display") // set your flag name and defaults however you wish
)

func main(){
    flag.Parse()
    profiler.SetProfileFlags(*profileLevelFlag) // this can be hard coded or sent in via command line.
}

Usage inside app: 
package myapp

import "github.com/zeroepix/goprofiler"


func WorkOn(x int){
    cp := profiler.New("WorkOn", 1)     // this is going to be a level 1 profiler
	defer cp.Finish()                   // important, this will perform the calculations and print it out. 
    //  You can call manually if you don't want to defer it.

    // do some initial setup
    cp.Tick("setup")                    // take a time measurement
    // do some work
    cp.Tick("work1")                    // take another time measurement
    // do some more work
    cp.Tick("work2")                    // take another measurement
    heavyWorkOn(x)                      // call another function
    cp.Tick("heavywork")                // this will measure the time between "work2" and the return from heavyWorkOn(x)
    // start some really tricky work, open a deep profiler here for this bit of code
    cp3 := profiler.New("hard work", 3) // this is flagged as a level 3
    // do a little bit of hard work
    cp3.Tick("hardwork1")
    // do more hard work
    cp3.Tick("hardwork2")
    // do even more hard work
    cp3.Tick("hardwork3")
    cp3.Finish()                        // finalise results of the hard work profiler
    cp.Tick("Hard Work")                // this will measure the hardwork section from the "heavywork" tick above to now
    return                              // no need to finalize cp as it was deferred
    
}

func heavyWorkOn(y int){
    cp := profiler.New("heavyWorkOn", 2) // this profiler will only show when level 2 is requested
    defer cp.Finish()
    // do some heavy work
    cp.Tick("heavywork1")
    // do some more
    cp.Tick("heavywork2")
    ...
    return
}