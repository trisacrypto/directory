import { Avatar, Box, HStack, Text } from '@chakra-ui/react';
import { Organization } from 'modules/dashboard/organization/organizationType';
import { BsChevronRight } from 'react-icons/bs';
import { Link } from 'react-router-dom';
import { APP_PATH } from 'utils/constants';
type AccountProps = {
  src?: string;
  alt?: string;
  onClose: () => void;
} & Partial<Organization>;

export const Account = ({ domain, name, src, id, onClose }: AccountProps) => {
  const orgLink = `${APP_PATH.SWITCH_ORGANIZATION}/${id}`;
  const selectOrgHandler = () => {
    onClose();
    window.location.href = orgLink;
  };
  return (
    <>
      <Link onClick={selectOrgHandler} to={orgLink} role="group">
        <HStack
          direction="row"
          w="100%"
          _groupHover={{ bg: '#BEE3F81d' }}
          padding={3}
          transition="background 400ms">
          <HStack width="100%" alignItems="center">
            <HStack>
              <Avatar bg="blue" name={name} src={src} fontWeight={700} />
              <Box>
                <Text fontWeight={700} data-testid="vaspName">
                  {name}
                </Text>
                <Text data-testid="vaspDomain">{domain}</Text>
              </Box>
            </HStack>
          </HStack>
          <BsChevronRight />
        </HStack>
      </Link>
    </>
  );
};
