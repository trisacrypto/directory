/* eslint-disable @typescript-eslint/no-unused-vars */
import React, { useEffect } from 'react';
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
  Button,
  useToast
} from '@chakra-ui/react';
import Loader from 'components/Loader';
import MemberModalContent from './MemberModalContent';
import { useFetchMember } from '../../hooks/useFetchMember';
import { Trans, t } from '@lingui/macro';
import { useSelector } from 'react-redux';
import { memberSelector } from '../../member.slice';
import Copy from './Copy';
interface MemberModalProps {
  isOpen: boolean;
  onClose: () => void;
  member: any;
}
const MemberModal = ({ isOpen, onClose, member: memberId }: MemberModalProps) => {
  const network = useSelector(memberSelector).members.network;
  const { member, isFetchingMember, error, wasMemberFetched } = useFetchMember({
    vaspId: memberId,
    network
  });
  const toast = useToast();

  useEffect(() => {
    if (error && error?.response?.status !== 451) {
      onClose();
      toast({
        description: error?.response?.data?.error,
        status: 'error',
        duration: 5000,
        position: 'top-right',
        isClosable: true
      });
    }
  }, [error, toast, onClose]);

  return (
    <Flex>
      <Box w="full">
        <Modal
          closeOnOverlayClick={false}
          isOpen={isOpen}
          onClose={onClose}
          size="lg"
          data-testid="member-modal">
          <ModalOverlay />
          <ModalContent width={'100%'} maxHeight={'1000px'}>
            {wasMemberFetched && !isFetchingMember && !error && (
              <>
                <ModalHeader data-testid="confirmation-modal-header" textAlign={'center'}>
                  {member?.data?.summary?.name}
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
                    <Copy data={member} />
                  </HStack>
                </ModalFooter>
              </>
            )}

            {isFetchingMember && (
              <ModalBody pb={6}>
                <Loader />
              </ModalBody>
            )}
          </ModalContent>
        </Modal>
      </Box>
    </Flex>
  );
};

export default MemberModal;
