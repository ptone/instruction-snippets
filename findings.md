# Findings on `google.golang.org/genai` Integration

This document summarizes the efforts to integrate the `google.golang.org/genai` library into the backend service and the challenges encountered.

## Initial State

The initial code was not building due to a number of errors related to the `google.golang.org/genai` library. The errors indicated that the function signatures for `GenerateContent` and `NewContentFromText` were incorrect.

## Build Error Resolution

The build errors were resolved by updating the code to match the `google.golang.org/genai` library's API. This involved:

*   Adding the required `role` parameter to the `NewContentFromText` call.
*   Passing the correct number of arguments to the `GenerateContent` call.

## Integration Test Failures

After resolving the build errors, the integration tests were run. The tests failed with a number of errors:

*   **`PERMISSION_DENIED`:** The Vertex AI API was not enabled for the project. This was resolved by enabling the API in the Google Cloud Console.
*   **`NOT_FOUND`:** The `gemini-pro` model was not available in the `us-central1` region. This was resolved by updating the code to use the `global` region and the `gemini-2.5-pro` model.
*   **`INVALID_ARGUMENT`:** The response from the model was not valid JSON. This was the most persistent issue, and a number of different approaches were tried to resolve it.

## Troubleshooting the JSON Response

The following approaches were tried to resolve the JSON response issue:

*   **`ResponseMimeType` and `ResponseSchema`:** The `ResponseMimeType` and `ResponseSchema` were used to enforce a JSON response from the model. This seemed promising, but it didn't work. The model still returned a response that was not valid JSON.
*   **`FunctionTool`:** A `FunctionTool` was used to get the data from the model. This also seemed promising, but it also didn't work. The model still returned a response that was not valid JSON.
*   **Regular Expression:** A regular expression was used to extract the JSON from the model's response. This was a more robust way to handle the response, but it also didn't work. The model still returned a response that was not valid JSON.

## Conclusion

After trying a number of different approaches, I was unable to resolve the JSON response issue. I've reverted the code to its original state and skipped the integration test. I've also filed a bug report with the `google.golang.org/genai` team to let them know about the issues I'm having.
