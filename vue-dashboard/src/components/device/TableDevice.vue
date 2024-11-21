<template>

    <div class="ttable__container">
 
        <input class="ttable__search mt-8" v-model="searchTerm" type="text" placeholder="Search...">

        <table class="ttable  mt-8">
            <thead>
                <tr>
                    <th class="cursor-pointer">
                        #
                        <svg class="t-sort-arrow t-sort-arrow--active">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Name
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Device ID
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Network 
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Firmware
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer ">
                        Latitude
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Longitude
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Occupied
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        White List
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Black List
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th>

                    </th>
                </tr>
            </thead>
            <tbody>

                <tr>
                    <td>1</td>
                    <td>Michael Scott</td>
                    <td>860226067572735</td>
                    <td>NB-IoT</td>
                    <td>3.2</td>
                    <td>35.928676</td>
                    <td>14.418638</td>
                    <td>true</td>
                    <td>false</td>
                    <td>false</td>
                    <td>
                        <div class="t-btns">
                            <a class="t-btns__btn">
                                <svg class="t-btns__icon">
                                    <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use>
                                </svg>
                            </a>
                        </div>
                    </td>
                </tr>

                <tr>
                    <td>2</td>
                    <td>
                        <input class="ttable__input" type="text" value="Pam Beesly">
                    </td>
                    <td>860226067572735</td>
                    <td>NB-IoT</td>
                    <td>3.2</td>
                    <td>
                        <input class="ttable__input w-24" type="text" value="35.928676">
                    </td>
                    <td>
                        <input class="ttable__input w-24" type="text" value="14.418638">
                    </td>
                    <td>true</td>
                    <td>
                        <select>
                            <option value="active">false</option>
                            <option value="deactive">true</option>                  
                        </select>
                    </td>
                    <td>
                        <select>
                            <option value="active">false</option>
                            <option value="deactive">true</option>                 
                        </select>
                    </td>
                
                    <td>
                        <div class="t-btns">
                            <cite class="t-btns__cite">Update Device</cite>
                            <a class="t-btns__btn t-btns__btn--yes">
                                Yes
                            </a>
                            <a class="t-btns__btn t-btns__btn--no">
                                No
                            </a>
                        </div>
                    </td>
                </tr>

                <tr>
                    <td>3</td>
                    <td>Jim Halpert</td>
                    <td>860226067572735</td>
                    <td>NB-IoT</td>
                    <td>3.2</td>
                    <td>35.928676</td>
                    <td>14.418638</td>
                    <td>true</td>
                    <td>false</td>
                    <td>false</td>
                    <td>
                        <div class="t-btns">
                            <a class="t-btns__btn">
                                <svg class="t-btns__icon">
                                    <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use>
                                </svg>
                            </a>
                        </div>
                    </td>
                </tr>
                <tr>
                    <td>4</td>
                    <td>Dwight Schrute</td>
                    <td>860226067572735</td>
                    <td>NB-IoT</td>
                    <td>3.2</td>
                    <td>35.928676</td>
                    <td>14.418638</td>
                    <td>true</td>
                    <td>false</td>
                    <td>false</td>
                    <td>
                        <div class="t-btns">
                            <a class="t-btns__btn">
                                <svg class="t-btns__icon">
                                    <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use>
                                </svg>
                            </a>
                        </div>
                    </td>
                </tr>

                <tr>
                    <td>5</td>
                    <td>Angela Martin</td>
                    <td>860226067572735</td>
                    <td>NB-IoT</td>
                    <td>3.2</td>
                    <td>35.928676</td>
                    <td>14.418638</td>
                    <td>true</td>
                    <td>false</td>
                    <td>false</td>
                    <td>
                        <div class="t-btns">
                            <a class="t-btns__btn">
                                <svg class="t-btns__icon">
                                    <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use>
                                </svg>
                            </a>
                        </div>
                    </td>
                </tr>

            </tbody>
        </table>

    </div>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { onMounted, ref } from 'vue';

// -- Data -------------------------------------------------------------
const searchTerm = ref("");


</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped> 

// Import global SCSS styles for consistent styling across the application


</style>
