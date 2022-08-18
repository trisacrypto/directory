import React from 'react';
import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogContent,
  AlertDialogOverlay,
  Button,
  useDisclosure,
  AlertDialogCloseButton
} from '@chakra-ui/react';
import { Trans } from '@lingui/react';

interface ModalProps {}
const ModalAlert = (props: any) => {
  const cancelRef: any = React.useRef();

  return (
    <>
      <AlertDialog
        motionPreset="slideInBottom"
        leastDestructiveRef={cancelRef}
        onClose={props.onClose}
        isOpen={props.isOpen}
        closeOnOverlayClick={false}>
        <AlertDialogOverlay />

        <AlertDialogContent>
          <AlertDialogHeader>{props.header}</AlertDialogHeader>
          <AlertDialogBody>{props.message}</AlertDialogBody>
          <AlertDialogFooter>
            <Button ref={cancelRef} onClick={props.onClose}>
              <Trans id="No">No</Trans>
            </Button>
            <Button colorScheme="green" ml={3} onClick={props.handleYesBtn}>
              <Trans id="Yes">Yes</Trans>
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
};

export default ModalAlert;
