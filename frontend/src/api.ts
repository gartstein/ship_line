// Set the API base URL from an environment variable.
// When using Create React App, env vars must be prefixed with REACT_APP_.
// Fallback to "http://localhost:8080/v1" if the env variable is not set.
const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || "http://localhost:8080/v1";

export const fetchPackSizes = async () => {
    const response = await fetch(`${API_BASE_URL}/pack-sizes`);
    if (!response.ok) throw new Error(`Error ${response.status}: Failed to fetch pack sizes.`);

    const data = await response.json();
    return { pack_sizes: Array.isArray(data.pack_sizes) ? data.pack_sizes : [] }; // Ensure it's an array
};

export const updatePackSizes = async (packSizes: number[]) => {
    const response = await fetch(`${API_BASE_URL}/pack-sizes`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ pack_sizes: packSizes }),
    });

    if (!response.ok) throw new Error(`Error ${response.status}: Failed to update pack sizes.`);
    return response.json();
};

export const calculatePacks = async (items: number) => {
    const response = await fetch(`${API_BASE_URL}/calc?items=${items}`);
    if (!response.ok) {
        throw new Error(`Error ${response.status}: Failed to calculate pack sizes.`);
    }
    const data = await response.json();
    const formattedResult = Object.entries(data.packsUsed).map(([pack, count]) => ({
        pack: parseInt(pack, 10),
        count: count as number,
    }));

    return {
        itemsOrdered: data.itemsOrdered,
        totalItemsUsed: data.totalItemsUsed,
        result: formattedResult,
    };
};

export const deletePackSize = async (size: number) => {
    const response = await fetch(`${API_BASE_URL}/pack-sizes/${size}`, {
        method: "DELETE",
    });
    if (!response.ok) {
        throw new Error("Failed to delete pack size");
    }
    return response.json();
};