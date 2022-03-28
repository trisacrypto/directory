import { Select } from '@chakra-ui/react';

const LanguagesDropdown: React.FC = () => {
  return (
    <Select w="100%" maxW="100">
      <option value="option1">🇬🇧</option>
      <option value="option2">🇫🇷</option>
    </Select>
  );
};

export default LanguagesDropdown;
