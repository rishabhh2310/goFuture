package main
import(
"fmt"
"context"
"time"
)

type Future interface {
  get() Result
  getWithTimeout(duration time.Duration) Result
  isComplete() bool
  isCancelled() bool 
  cancel()
}


type Result struct{
  resultValue interface{}
  error error
}





type FutureTask struct{
  success bool   
  done bool      
  error error    
  result Result 
  interfaceChannel <-chan Result 
}



func (futureTask *FutureTask) get() Result{
  if(futureTask.done){
    return futureTask.result
  }
  ctx := context.Background()
  return futureTask.getWithContext(ctx)
}




func (futureTask *FutureTask) getWithTimeout(timeout time.Duration) Result{
  if(futureTask.done){
    return futureTask.result
  }
  ctx, cancel := context.WithTimeout(context.Background(), timeout)
  defer cancel()
  return futureTask.getWithContext(ctx)
}



func (futureTask *FutureTask) getWithContext(ctx context.Context) Result{
  select {
  case <-ctx.Done():
    futureTask.done = true
    futureTask.success = false
    futureTask.error = &TimeoutError{errorString:"Request Timeout!"}
    futureTask.result = Result{resultValue:nil,error:futureTask.error}
    return futureTask.result
  case futureTask.result = <-futureTask.interfaceChannel:
    if(futureTask.result.error!=nil){
      futureTask.done = true
      futureTask.success = false
      futureTask.error = futureTask.result.error
    }else{
      futureTask.success = true
      futureTask.done = true
      futureTask.error = nil
    }
    return futureTask.result
  }
}



func (futureTask *FutureTask) isComplete() bool{
  if(futureTask.done){
    return true
  }else{
    return false;
  }
}



func (futureTask *FutureTask) isCancelled() bool{
  if(futureTask.done){
    if(futureTask.error!=nil && futureTask.error.Error()=="Cancelled Manually"){
      return true;
    }
  }
  return false;
}


func (futureTask *FutureTask) cancel(){
  if(futureTask.isComplete() || futureTask.isCancelled()){
    return;
  }
  interruptionError := &InterruptError{errorString:"Cancelled Manually"}
  futureTask.done = true
  futureTask.success = false
  futureTask.error = interruptionError
  futureTask.result =Result{resultValue:nil,error:interruptionError}
}



func ReturnAFuture(task func() (Result)) *FutureTask{
channelForExecution := make(chan Result)
futureObject := FutureTask{
  success :          false,
  done    :          false,
  error   :          nil,
  result  :      Result{},
  interfaceChannel : channelForExecution,
 }
go func(){
  defer close(channelForExecution)
  resultObject := task()
  channelForExecution <- resultObject
 }()
 return &futureObject
}


func main(){
futureInstance1:= ReturnAFuture(func() (Result){
   var res interface{}
   res=30+23
   time.Sleep(2*time.Second)
   return Result{resultValue:res}
 })
 futureInstance2:= ReturnAFuture(func() (Result){
   var res interface{}
   res="40"
   return Result{resultValue:res}
 })
 fmt.Println("before future")
 fmt.Println(futureInstance2.getWithTimeout(time.Second))
 fmt.Println(futureInstance1.get())
}




func (e *TimeoutError) Error() string{
  return e.errorString
}
func (e *InterruptError) Error() string{
  return e.errorString
}
type TimeoutError struct{
  errorString string
}
type InterruptError struct{
  errorString string
}
