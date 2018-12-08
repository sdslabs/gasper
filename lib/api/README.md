# `github.com/sdslabs/SWS/lib/api`

This package consists of methods that have a direct affect on the user requests. Some guidelines that need to be followed:

- For every exportable function, use the type `ResponseError` as in package `.../SWS/lib/types`.

  ```go
  // The NewResponseError function takes 3 arguments
  // Status Code: int
  // Message: string -> Any customized message for the error
  // Error: error -> If any error is thrown and response message is the same

  // When message is customized
  err := types.NewResponseError(400, "Invalid user input", nil)

  // When message is same as someother error message
  err1 := someErrorThrowingFunc()
  err2 := types.NewResponseError(500, "", err1)
  ```

  \* Customized error messages are required as many times sensitive information is revealed while sending the response back

  This will be further used as follows:

  ```go
  func controllerFunc(c *gin.Context) {
      // Some code here
      // If an error is generated
      c.JSON(err.Status(), gin.H{
          "error": err.Reason(),
      })
  }
  ```
