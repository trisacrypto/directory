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
import { Trans } from '@lingui/macro';
import { useSelector } from 'react-redux';
import { memberSelector } from '../../member.slice';
interface MemberModalProps {
  isOpen: boolean;
  onClose: () => void;
  member: any;
}
const MemberModal = ({ isOpen, onClose, member: memberId }: MemberModalProps) => {
  const network = useSelector(memberSelector).members.network;
  const { member, isFetchingMember, error } = useFetchMember({ vaspId: memberId, network });
  if (error) {
    // close the modal for now , later we need to show a toast message
    onClose();
  }
  return (
    <Flex>
      <Box w="full">
        <Modal
          closeOnOverlayClick={false}
          isOpen={isOpen}
          onClose={onClose}
          data-testid="member-modal">
          <ModalOverlay />
          <ModalContent width={'100%'}>
            {isFetchingMember ? (
              <Loader />
            ) : (
              <>
                <ModalHeader data-testid="confirmation-modal-header" textAlign={'center'}>
                  {member?.summary?.name}
                </ModalHeader>
                <ModalCloseButton data-testid="close-btn-icon" />

                <ModalBody pb={6}>
                  <MemberModalContent member={member} />
                </ModalBody>

                <ModalFooter>
                  <HStack width="100%" justifyContent="center" alignItems="center">
                    <Button bg={'black'} onClick={onClose} data-testid="modal-close-button">
                      <Trans>Close</Trans>
                    </Button>
                    <Button bg={'#FF7A59'} color={'white'}>
                      <Trans>Copy</Trans>
                    </Button>
                  </HStack>
                </ModalFooter>
              </>
            )}
          </ModalContent>
        </Modal>
      </Box>
    </Flex>
  );
};

export default MemberModal;
