/* eslint-disable @typescript-eslint/no-unused-vars */
import React from 'react';
import {
  Box,
  Flex,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Modal,
  HStack,
  Button
} from '@chakra-ui/react';
import Loader from 'components/Loader';
import MemberModalContent from './MemberModalContent';
import { useFetchMember } from '../../hooks/useFetchMember';
interface MemberModalProps {
  isOpen: boolean;
  onClose: () => void;
  member: any;
}
const MemberModal = ({ isOpen, onClose, member: memberId }: MemberModalProps) => {
  const { member, isFetchingMember } = useFetchMember(memberId);
  return (
    <>
      <Flex>
        <Box w="full">
          {isFetchingMember && <Loader />}
          {member && (
            <Modal closeOnOverlayClick={false} isOpen={isOpen} onClose={onClose}>
              <ModalOverlay />
              <ModalContent width={'100%'}>
                <ModalHeader data-testid="confirmation-modal-header" textAlign={'center'}>
                  {member?.name}
                </ModalHeader>
                <ModalCloseButton />

                <ModalBody pb={6}>
                  <MemberModalContent member={member} />
                </ModalBody>

                <ModalFooter>
                  <HStack width="100%" justifyContent="center" alignItems="center">
                    <Button onClick={onClose}>Close</Button>
                    <Button>Copy</Button>
                  </HStack>
                </ModalFooter>
              </ModalContent>
            </Modal>
          )}
        </Box>
      </Flex>
    </>
  );
};

export default MemberModal;
