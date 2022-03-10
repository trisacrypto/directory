import { Meta, Story } from "@storybook/react";
import DetailsCard from ".";

interface DetailsCardProps {
  type: string;
  data: any;
}
export default {
  title: "components/DetailsCard",
  component: DetailsCard,
} as Meta;

const Template: Story = (args: DetailsCardProps) => <DetailsCard {...args} />;

export const Default = Template.bind({});
Default.args = {};
