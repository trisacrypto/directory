import React from 'react';
import { Box, Alert, AlertDescription } from '@chakra-ui/react';
type TError = {
  message: string;
  handleClose?: () => void;
};

const SuccessMessage: React.FC<TError> = ({ message }) => {
  return (
    <Box my={4}>
      <Alert status="success" borderRadius={4} data-testid="success__alert">
        <AlertDescription>{message}</AlertDescription>
        {/* <CloseButton position="absolute" right="8px" top="8px" onClick={handleClose} /> */}
      </Alert>
    </Box>
  );
};
export default SuccessMessage;
