// src/components/PackSizeForm.tsx
import React, { useState } from "react";
import { TextField, Button } from "@mui/material";
import { updatePackSizes, fetchPackSizes } from "../api";

const PackSizeForm = ({ onPackSizeUpdate }: { onPackSizeUpdate: (newSizes: number[]) => void }) => {
    const [newPackSizes, setNewPackSizes] = useState("");
    const [error, setError] = useState("");

    const handleUpdate = async () => {
        if (!newPackSizes.trim()) {
            setError("Please enter valid numbers separated by commas.");
            return;
        }

        // Convert comma-separated input into an array of integers
        const sizesArray = newPackSizes
            .split(",")
            .map((size) => parseInt(size.trim(), 10))
            .filter((size) => !isNaN(size)); // Remove invalid numbers

        if (sizesArray.length === 0) {
            setError("Please enter at least one valid number.");
            return;
        }

        try {
            //Fetch existing sizes and append new ones
            const existingSizes = await fetchPackSizes();
            const updatedSizes = [...existingSizes.pack_sizes, ...sizesArray];

            await updatePackSizes(updatedSizes);
            setNewPackSizes("");

            onPackSizeUpdate(updatedSizes);
        } catch (err) {
            setError((err as Error).message);
        }
    };

    return (
        <div className="pack-form">
            <TextField
                fullWidth
                label="New Pack Sizes (comma-separated)"
                variant="outlined"
                value={newPackSizes}
                onChange={(e) => setNewPackSizes(e.target.value)}
                className="input-field"
            />
            <Button variant="contained" color="primary" onClick={handleUpdate} className="submit-btn">
                Add Pack Sizes
            </Button>
            {error && <p className="error">{error}</p>}
        </div>
    );
};

export default PackSizeForm;