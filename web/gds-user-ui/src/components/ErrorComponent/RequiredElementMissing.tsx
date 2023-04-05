import { Alert, AlertTitle, AlertDescription, Stack } from '@chakra-ui/react';
interface RequiredElementMissingProps {
  elementKey?: number;
}

const RequiredElementMissing = ({ elementKey }: RequiredElementMissingProps) => {
  console.log('elementKey', elementKey);

  return (
    <Alert status="error" borderRadius="lg" my={4}>
      <Stack>
        <AlertTitle>Required element(s) missing </AlertTitle>
        <AlertDescription>
          Please make sure you have filled out all required fields marked with an asterisk (*).
        </AlertDescription>
      </Stack>
    </Alert>
  );
};

export default RequiredElementMissing;
