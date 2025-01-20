<script setup lang="ts">
import { ref, onMounted } from 'vue';
import QrcodeVue from 'qrcode.vue'

function Uuid(length: number, radix: number): string {
  const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"; // UUID Seed from EAS
  const charArray = chars.split("");
  let uuid: string[] = [];

  // Initialize the random number generator with a seed based on the current time
  const rand = () => Math.floor(Math.random() * radix);

  if (length > 0) {
    // Compact form
    for (let i = 0; i < length; i++) {
      const index = rand();
      uuid.push(charArray[index]);
    }
  } else {
    // RFC 4122, version 4 form
    uuid = new Array(36).fill("");
    uuid[8] = uuid[13] = uuid[18] = uuid[23] = "-";
    uuid[14] = "4";

    for (let i = 0; i < 36; i++) {
      if (uuid[i] === "") {
        const r = rand();
        if (i === 19) {
          uuid[i] = charArray[(r & 0x3) | 0x8]; // Set to RFC 4122 High Bit
        } else {
          uuid[i] = charArray[r];
        }
      }
    }
  }

  return "smdljwxt" + uuid.join(""); // Magic string from the web side
}

const uuid16 = ref<string>();
const jsonData = ref<string>();
const ifScaned = ref<boolean>(false);
const jsonObj = ref<null>(null);

const fetchData = async () => {
  try {
    // Construct request URL
    const url = `/api/getmark?uuid16=${uuid16.value}`;
    const response = await fetch(url);
    if (response.ok) {
      const data = await response.json();
      jsonData.value = JSON.stringify(data);
      ifScaned.value = true;
    } else {
      console.error('请求失败:', response.status);
    }
  } catch (error) {
    console.error('请求出错:', error);
  }
};

onMounted(
  () => {
    uuid16.value = Uuid(16, 32)
    fetchData().then(() => {
    if (jsonData.value) {
      jsonObj.value = JSON.parse(jsonData.value);
    }
  });
  }
);

</script>

<template>
  <div class="mx-3">
    <p class="font-bold text-center text-3xl my-5 mx-3">
      请使用喜鹊扫描下方二维码以完成成绩收集
    </p>
    <p class="font-normal text-center text-2xl mx-3">注意: 长按该二维码可保存图片</p>
  </div>
  <div>
    <p class="font-normal text-center text-1xl my-5">本次您的登录 UUID 是 : </p>
    <p class="font-mono text-center text-1xl my-5">{{ uuid16 }}</p>
  </div>
  <div>
    <qrcode-vue :value="uuid16" :size="200" level="H" class="place-self-center" />
  </div>
  <div class="my-3">
    <p class="text-center m-3">扫码后稍等片刻，下方会显示您的信息</p>
    <p class="text-center m-3">请确认您的成绩数据和个人信息是否正确，如有误，请截取本界面并联系 QQ2954582482 ：</p>
  </div>
  <p class="text-center m-3" v-if="jsonObj">您正在成功！请坐和放宽！</p>
  <div v-if="jsonObj" class="grid grid-flow-row place-items-center m-4">
    <div>
      <p>学号: {{ jsonObj["ID"] }}</p>
      <p>姓名: {{ jsonObj["Name"] }}</p>
      <p>班级: {{ jsonObj["ClassName"] }}</p>
      <p>学院: {{ jsonObj["Institution"] }}</p>
    </div>
    <table class="table-auto border-spacing-2 border border-slate-500">
      <thead>
        <tr>
          <th class="border border-slate-600 text-left p-3">科目</th>
          <th class="border border-slate-600 text-center p-3">成绩</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(value, key) in jsonObj['Scores']" :key="key">
          <th class="border border-slate-600 text-left p-3"> {{ key }} </th>
          <th class="border border-slate-600 text-center"> {{ value }} </th>
        </tr>
      </tbody>
    </table>
  </div>
</template>