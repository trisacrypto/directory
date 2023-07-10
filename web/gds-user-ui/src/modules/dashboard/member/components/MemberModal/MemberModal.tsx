/* eslint-disable @typescript-eslint/no-unused-vars */
import React, { useState } from 'react';
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
} from '@chakra-ui/react';
import Loader from 'components/Loader';
import MemberModalContent from './MemberModalContent';
import { useFetchMember } from '../../hooks/useFetchMember';
import { Trans } from '@lingui/macro';
import { memberDetailMock } from '../../__mocks__';
interface MemberModalProps {
  isOpen: boolean;
  onClose: () => void;
  member: any;
}
const MemberModal = ({ isOpen, onClose, member: memberId }: MemberModalProps) => {
  const { /* member, */ /* isFetchingMember */ } = useFetchMember(memberId);
  const mock = memberDetailMock;


  return (
    <>
      <Flex>
        <Box w="full">
          {/* {isFetchingMember && <Loader />} */}
          {mock && (
            <Modal
              closeOnOverlayClick={false}
              isOpen={isOpen}
              onClose={onClose}
              data-testid="member-modal">
              <ModalOverlay />
              <ModalContent width={'100%'}>
                <ModalHeader data-testid="confirmation-modal-header" textAlign={'center'}>
                  {mock?.summary?.name}
                </ModalHeader>
                <ModalCloseButton data-testid="close-btn-icon" />

                <ModalBody pb={6}>
                  <MemberModalContent member={mock} />
                </ModalBody>

                <ModalFooter>
                  <HStack width="100%" justifyContent="center" alignItems="center">
                    <Button bg={'black'} onClick={onClose} data-testid="modal-close-button">
                      <Trans>Close</Trans>
                    </Button>
                    <Button bg={'#FF7A59'} color={'white'}>
                      <Trans>Export</Trans>
                    </Button>
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
