import React, { useEffect, useState } from 'react';
import MemberModal from './MemberModal';
import { useDisclosure, Button, HStack } from '@chakra-ui/react';
import { BsEye } from 'react-icons/bs';
import { useFetchMember } from '../../hooks/useFetchMember';
interface ShowMemberModalProps {
  memberId: any;
}
const ShowMemberModal: React.FC<ShowMemberModalProps> = ({ memberId }) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [shouldOpenModal, setShouldOpenModal] = useState(false);

  const handleOpenModal = () => {
    setShouldOpenModal(true);
  };

  const handleCloseModal = () => {
    setShouldOpenModal(false);
    onClose();
  };

  useEffect(() => {
    if (shouldOpenModal) {
      onOpen();
    }
  }, [shouldOpenModal, onOpen]);

  return (
    <>
      <HStack width="100%" justifyContent="center" alignItems="center">
        <Button
          data-testid="member-modal-button"
          onClick={handleOpenModal}
          color="blue"
          bg={'transparent'}
          _hover={{
            bg: 'transparent'
          }}
          _focus={{
            bg: 'transparent'
          }}>
          <BsEye fontSize="24px" />
        </Button>
      </HStack>
      {shouldOpenModal && (<MemberModal isOpen={isOpen} onClose={handleCloseModal} member={memberId} />)}
    </>
  );
};

export default ShowMemberModal;
