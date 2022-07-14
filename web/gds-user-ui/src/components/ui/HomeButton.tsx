import { Button, Icon, ButtonProps, useColorModeValue } from '@chakra-ui/react';
import { AiOutlineHome } from 'react-icons/ai';

type HomeButtonProps = {
  link: string;
};

const HomeButton: React.FC<HomeButtonProps & ButtonProps> = ({ link, ...props }) => {
  const textColor = useColorModeValue('gray.800', '#F7F8FC');

  return (
    <Button
      as={'a'}
      role="group"
      href={link}
      {...props}
      backgroundColor={'transparent'}
      _hover={{ background: 'black', color: 'white' }}
      _focus={{ background: 'red.100', color: 'red.500' }}>
      <Icon as={AiOutlineHome} fontSize={'24'} color={textColor} _groupHover={{ color: 'white' }} />
    </Button>
  );
};

export default HomeButton;
