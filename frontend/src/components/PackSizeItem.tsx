// src/components/PackSizeItem.tsx
import React from "react";
import { deletePackSize } from "../api";

interface PackSizeItemProps {
    size: number;
    onDelete: (size: number) => void;
}

const PackSizeItem: React.FC<PackSizeItemProps> = ({ size, onDelete }) => {
    const handleClick = async () => {
        if (window.confirm(`Do you really want to remove size ${size}?`)) {
            try {
                await deletePackSize(size);
                onDelete(size);
            } catch (error) {
                console.error("Failed to delete pack size:", error);
            }
        }
    };

    return (
        <div className="pack-size-item" onClick={handleClick}>
            {size}
        </div>
    );
};

export default PackSizeItem;