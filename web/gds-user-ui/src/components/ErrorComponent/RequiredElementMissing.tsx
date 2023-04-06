import { Alert, AlertIcon } from '@chakra-ui/react';
interface RequiredElementMissingProps {
  elementKey?: number;
}

const RequiredElementMissing = ({ elementKey }: RequiredElementMissingProps) => {
  console.log('k', elementKey);

  return (
    <Alert status="error" borderRadius="lg" my={4}>
      <AlertIcon />
      Please make sure you have filled out all required fields marked with an asterisk (*).
    </Alert>
  );
};

export default RequiredElementMissing;
