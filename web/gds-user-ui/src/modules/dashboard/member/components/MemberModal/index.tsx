import React from 'react';
import MemberModal from './MemberModal';
import { useDisclosure, Button, HStack } from '@chakra-ui/react';
import { BsEye } from 'react-icons/bs';
import { useFetchMember } from '../../hooks/useFetchMember';
interface ShowMemberModalProps {
  memberId: any;
}
const ShowMemberModal: React.FC<ShowMemberModalProps> = ({ memberId }) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const { isFetchingMember } = useFetchMember(memberId);

  return (
    <>
      <HStack width="100%" justifyContent="center" alignItems="center">
        <Button
          isLoading={isFetchingMember}       
          data-testid="member-modal-button"
          onClick={onOpen}
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
      <MemberModal isOpen={isOpen} onClose={onClose} member={memberId} />
    </>
  );
};

export default ShowMemberModal;
