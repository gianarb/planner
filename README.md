## planner

Planner is a library I wrote to implement reactive planning in Go.

First what is reactive planning? I will leave you a few links:

* [Reactive planning from Wikipedia](https://en.wikipedia.org/wiki/Reactive_planning)
* [The book "Thinking in Systems"](https://gianarb.it/blog/thinking-in-systems-donella-meadows-review)
* [AWS re:Invent 2018: Close Loops & Opening Minds: How to Take Control of Systems, Big & Small ARC337](https://www.youtube.com/watch?v=O8xLxNje30M&feature=emb_title)
* [Me explaining on Twitch this library with an example](https://www.twitch.tv/videos/770165588)

## Reactive planning implementation in planner

In short a reactive plan empower a humanoid robot the reconcile and stay up when
you push it.

Or it keeps your Kubernetes resources such as Pods, Services, Ingress and
Deployment up and running.

It is made of three parts:

1. A plan
2. A set of procedures
3. A scheduler

A **procedure** is the smallest unit of work that you can think about:

1. Create an AWS EC2
2. Control the location of a device
3. Do a server healthcheck

A **plan** is a set of **procedures**. If we look at how a replication
controller works (AWS Autoscaling Group or Kubernetes Replicaset) we can think
about a set of common steps like:

1. Find the list of available servers or pods for that particular replicaset or
   autoscaling group
2. Do a healthcheck on all the servers/pod
3. Based on how many of them are healthy there are other 2 steps:
    4. Create a new EC2/Pod
    5. Delete an EC2/Pod

Those are 5 **procedures**, together they are a Reconciliation plan.

The **scheduler** takes a plan and execute it until there is nothing left to do.

A plan has a `Create` function that calculates and returns the **procedures**
that has to be executed. Iteration over iteration the number of **procedures** can
change.

Your **plan** can succeed at the first iteration, it means that the second one will
return zero **procedures** and the scheduler will stop executing the plan.

If there are left over action to be done they will be picked up during a future
execution.

A scheduler stops to execute a plan only if:

1. There are not left over procedures for a particular plan (it is all done!
   Great)
2. A procedure returns an `error`

Each **procedure** can return multiple procedures, in this way you mitigate the
amount of errors you have to return, zero is your target! Any `error` you
encounter is an opportunity to code a mitigation as separate procedure. Sometime
you have to just wait, sometime you can trigger a page and wait until human
fixes it.

## Example

I would like to measure random and luck! So I would like to write a program that
given a number tries its best to get there just via random additions and
subtraction [you can start from here](https://play.golang.com/p/0LuIoMtp10f)
using reactive planning.

If you execute the program in its current form it will get this output:

```console
1.257894e+09	info	planner@v0.0.1/scheduer.go:41	Started execution plan count_plan	{"execution_id": "befecb26-1c94-4a61-8305-c9b40aa63331"}
1.257894e+09	info	planner@v0.0.1/scheduer.go:59	Plan executed without errors.	{"execution_id": "befecb26-1c94-4a61-8305-c9b40aa63331", "execution_time": "0s", "step_executed": 20}
```

The plan is a success in 20 steps. Because as you can see there is only one
procedures who increments by one. So starting from zero it takes 20 steps to get
to my desired number 20. That's great to get the vibe of the project, but we
need more.

Suggested evolution:

* Change the AddNumber to use a randomly generated number
* Add a condition in the `Create` that looks like this:

```go
	if p.current < p.Target {
		return []planner.Procedure{&AddNumber{plan: p}}, nil
	} else {
		return []planner.Procedure{&SubtractNumber{plan: p}}, nil
    }
```
And write `SubtractNumber` in the same way `AddNumber` works but with `-`.

## Are you using this library

Add your project to [ADOPTERS.md](./ADOPTERS.md)
