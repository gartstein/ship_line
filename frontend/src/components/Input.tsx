import React from "react";
import { TextField } from "@mui/material";

interface InputProps {
    value: string;
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

const Input: React.FC<InputProps> = ({ value, onChange }) => {
    return (
        <TextField
            label="Number of Items"
            variant="outlined"
            fullWidth
            value={value}
            onChange={onChange}
        />
    );
};

export default Input;