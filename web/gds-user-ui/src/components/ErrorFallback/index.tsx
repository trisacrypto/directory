import { Button, Code, Text } from '@chakra-ui/react';

function ErrorFallback({ error, resetErrorBoundary }: any) {
  return (
    <div role="alert">
      <Text>Something went wrong:</Text>
      <Code>{error.message}</Code>
      <Button onClick={resetErrorBoundary}>Try again</Button>
    </div>
  );
}

export default ErrorFallback;
